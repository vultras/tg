package main

import (
	"log"
	"os"
	"strings"

	"github.com/mojosa-software/got/tg"
)

type BotData struct {
	Name string
}

type SessionData struct {
	Counter int
}

type MutateMessageWidget struct {
	Mutate func(string) string
}

func NewMutateMessageWidget(fn func(string) string) *MutateMessageWidget {
	ret := &MutateMessageWidget{}
	ret.Mutate = fn
	return ret
}

func (w *MutateMessageWidget) Serve(c *tg.Context, updates chan *tg.Update) error {
	for _, arg := range c.Args {
		c.Sendf("%v", arg)
	}
	for u := range updates {
		if u.Message == nil {
			continue
		}
		text := u.Message.Text
		c.Sendf("%s", w.Mutate(text))
	}
	return nil
}

func ExtractSessionData(c *tg.Context) *SessionData {
	return c.Session.Data.(*SessionData)
}

var (
	startScreenButton = tg.NewButton("🏠 To the start screen").
				ScreenChange("start")

	incDecKeyboard = tg.NewReply().Row(
		tg.NewButton("+").ActionFunc(func(c *tg.Context) {
			d := ExtractSessionData(c)
			d.Counter++
			c.Sendf("%d", d.Counter)
		}),
		tg.NewButton("-").ActionFunc(func(c *tg.Context) {
			d := ExtractSessionData(c)
			d.Counter--
			c.Sendf("%d", d.Counter)
		}),
	).Row(
		startScreenButton,
	)

	navKeyboard = tg.NewReply().
			WithOneTime(true).
			Row(
			tg.NewButton("Inc/Dec").ScreenChange("start/inc-dec"),
		).Row(
		tg.NewButton("Upper case").ActionFunc(func(c *tg.Context){
			c.ChangeScreen("start/upper-case", "this shit", "works")
		}),
		tg.NewButton("Lower case").ScreenChange("start/lower-case"),
	).Row(
		tg.NewButton("Send location").ScreenChange("start/send-location"),
	)

	sendLocationKeyboard = tg.NewReply().
				Row(
			tg.NewButton("Send location").
				WithSendLocation(true).
				ActionFunc(func(c *tg.Context) {
					var err error
					if c.Message.Location != nil {
						l := c.Message.Location
						_, err = c.Sendf(
							"Longitude: %f\n"+
								"Latitude: %f\n"+
								"Heading: %d"+
								"",
							l.Longitude,
							l.Latitude,
							l.Heading,
						)
					} else {
						_, err = c.Sendf("Somehow location was not sent")
					}
					if err != nil {
						c.Sendf("%q", err)
					}
				}),
		).Row(
		startScreenButton,
	)

	// The keyboard to return to the start screen.
	navToStartKeyboard = tg.NewReply().Row(
		startScreenButton,
	)
)

var beh = tg.NewBehaviour().
	WithInitFunc(func(c *tg.Context) {
		// The session initialization.
		c.Session.Data = &SessionData{}
	}).WithScreens(
		tg.NewScreen("start", tg.NewPage(
				"The bot started!",
			).WithInline(
				tg.NewInline().Row(
					tg.NewButton("GoT Github page").
						WithUrl("https://github.com/mojosa-software/got"),
				),
			).WithReply(
				navKeyboard,
			),
		),
		tg.NewScreen("start/inc-dec", tg.NewPage(
				"The screen shows how "+
					"user separated data works "+
					"by saving the counter for each of users "+
					"separately. ",
			).WithReply(
				incDecKeyboard,
			).ActionFunc(func(c *tg.Context) {
				// The function will be calleb before serving page.
				d := ExtractSessionData(c)
				c.Sendf("Current counter value = %d", d.Counter)
			}),
		),

		tg.NewScreen("start/upper-case", tg.NewPage(
				"Type text and the bot will send you the upper case version to you",
			).WithReply(
				navToStartKeyboard,
			).WithSub(
				NewMutateMessageWidget(strings.ToUpper),
			),
		),

		tg.NewScreen("start/lower-case", tg.NewPage(
				"Type text and the bot will send you the lower case version",
			).WithReply(
				navToStartKeyboard,
			).WithSub(
				NewMutateMessageWidget(strings.ToLower),
			),
		),

		tg.NewScreen("start/send-location", tg.NewPage(
				"Send your location and I will tell where you are!",
			).WithReply(
				sendLocationKeyboard,
			).WithInline(
				tg.NewInline().Row(
					tg.NewButton(
						"Check",
					).WithData(
						"check",
					).ActionFunc(func(c *tg.Context) {
							d := ExtractSessionData(c)
							c.Sendf("Counter = %d", d.Counter)
					}),
				),
			),
		),
	).WithCommands(
		tg.NewCommand("start").
			Desc("start or restart the bot or move to the start screen").
			ActionFunc(func(c *tg.Context){
				c.Sendf("Your username is %q", c.Message.From.UserName)
				c.ChangeScreen("start")
			}),
		tg.NewCommand("hello").
			Desc("sends the 'Hello, World!' message back").
			ActionFunc(func(c *tg.Context) {
				c.Sendf("Hello, World!")
			}),
		tg.NewCommand("read").
			Desc("reads a string and sends it back").
			WidgetFunc(func(c *tg.Context, updates chan *tg.Update) error {
				c.Sendf("Type text and I will send it back to you")
				for u := range updates {
					if u.Message == nil {
						continue
					}
					c.Sendf("You typed %q", u.Message.Text)
					break
				}
				c.Sendf("Done")
				return nil
			}),
		tg.NewCommand("image").
			Desc("sends a sample image").
			ActionFunc(func(c *tg.Context) {
				img := tg.NewFile("media/cat.jpg").Image().Caption("A cat!")
				c.Send(img)
			}),
		tg.NewCommand("botname").
			Desc("get the bot name").
			ActionFunc(func(c *tg.Context) {
				bd := c.Bot.Data.(*BotData)
				c.Sendf("My name is %q", bd.Name)
			}),
	)

var gBeh = tg.NewGroupBehaviour().
	InitFunc(func(c *tg.GC) {
	}).
	WithCommands(
		tg.NewGroupCommand("hello").ActionFunc(func(c *tg.GC) {
			c.Sendf("Hello, World!")
		}),
		tg.NewGroupCommand("mycounter").ActionFunc(func(c *tg.GC) {
			d := c.Session().Data.(*SessionData)
			c.Sendf("Your counter value is %d", d.Counter)
		}),
	)

func main() {
	token := os.Getenv("BOT_TOKEN")

	bot, err := tg.NewBot(token)
	if err != nil {
		log.Panic(err)
	}
	bot = bot.
		WithBehaviour(beh).
		WithGroupBehaviour(gBeh).
		Debug(true)

	bot.Data = &BotData{
		Name: "Jay",
	}

	log.Printf("Authorized on account %s", bot.Api.Self.UserName)
	err = bot.Run()
	if err != nil {
		panic(err)
	}
}
