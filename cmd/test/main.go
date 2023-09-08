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

type UserData struct {
	Counter int
}

var (
	startScreenButton = tg.NewButton("🏠 To the start screen").
				ScreenChange("start")

	incDecKeyboard = tg.NewInline().Row(
		tg.NewButton("+").ActionFunc(func(c *tg.Context) {
			d := c.Session.Value.(*UserData)
			d.Counter++
			c.Sendf("%d", d.Counter)
		}),
		tg.NewButton("-").ActionFunc(func(c *tg.Context) {
			d := c.Session.Value.(*UserData)
			d.Counter--
			c.Sendf("%d", d.Counter)
		}),
	).Row(
		startScreenButton,
	)

	navKeyboard = tg.NewReply().
			WithOneTime(true).
			Row(
			tg.NewButton("Inc/Dec").ScreenChange("inc/dec"),
		).Row(
		tg.NewButton("Upper case").ScreenChange("upper-case"),
		tg.NewButton("Lower case").ScreenChange("lower-case"),
	).Row(
		tg.NewButton("Send location").ScreenChange("send-location"),
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
						_, err = c.Sendf("Somehow wrong location was sent")
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
		c.Session.Value = &UserData{}

	}). // On any message update before the bot created session.
	WithPreStartFunc(func(c *tg.Context){
		c.Sendf("Please, use the /start command to start the bot")
	}).WithScreens(
	tg.NewScreen("start").
		WithText(
			"The bot started!"+
				" The bot is supposed to provide basic"+
				" understand of how the API works, so just"+
				" horse around a bit to guess everything out"+
				" by yourself!",
		).WithReply(navKeyboard).
		// The inline keyboard with link to GitHub page.
		WithInline(
			tg.NewInline().Row(
				tg.NewButton("GoT Github page").
					WithUrl("https://github.com/mojosa-software/got"),
			),
		),

	tg.NewScreen("inc/dec").
		WithText(
			"The screen shows how "+
				"user separated data works "+
				"by saving the counter for each of users "+
				"separately. ",
		).
		WithReply(&tg.ReplyKeyboard{Keyboard: incDecKeyboard.Keyboard}).
		// The function will be called when reaching the screen.
		ActionFunc(func(c *tg.Context) {
			d := c.Session.Value.(*UserData)
			c.Sendf("Current counter value = %d", d.Counter)
		}),

	tg.NewScreen("upper-case").
		WithText("Type text and the bot will send you the upper case version to you").
		WithReply(navToStartKeyboard).
		ActionFunc(mutateMessage(strings.ToUpper)),

	tg.NewScreen("lower-case").
		WithText("Type text and the bot will send you the lower case version").
		WithReply(navToStartKeyboard).
		ActionFunc(mutateMessage(strings.ToLower)),

	tg.NewScreen("send-location").
		WithText("Send your location and I will tell where you are!").
		WithReply(sendLocationKeyboard).
		WithInline(
			tg.NewInline().Row(
				tg.NewButton("Check").
					WithData("check").
					ActionFunc(func(a *tg.Context) {
						d := a.Session.Value.(*UserData)
						a.Sendf("Counter = %d", d.Counter)
					}),
			),
		),
).WithCommands(
	tg.NewCommand("start").
		Desc("start the bot").
		ActionFunc(func(c *tg.Context){
			c.ChangeScreen("start")
		}),
	tg.NewCommand("hello").
		Desc("sends the 'Hello, World!' message back").
		ActionFunc(func(c *tg.Context) {
			c.Sendf("Hello, World!")
		}),
	tg.NewCommand("read").
		Desc("reads a string and sends it back").
		ActionFunc(func(c *tg.Context) {
			c.Sendf("Type some text:")
			msg, err := c.ReadTextMessage()
			if err != nil {
				return
			}
			c.Sendf("You typed %q", msg)
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
			bd := c.Bot.Value().(*BotData)
			c.Sendf("My name is %q", bd.Name)
		}),
)

func mutateMessage(fn func(string) string) tg.ActionFunc {
	return func(c *tg.Context) {
		for {
			msg, err := c.ReadTextMessage()
			if err == tg.NotAvailableErr {
				break
			} else if err != nil {
				panic(err)
			}

			_, err = c.Sendf("%s", fn(msg))
			if err != nil {
				panic(err)
			}
		}
	}
}

var gBeh = tg.NewGroupBehaviour().
	InitFunc(func(c *tg.GC) {
	}).
	WithCommands(
		tg.NewGroupCommand("hello").ActionFunc(func(c *tg.GC) {
			c.Send("Hello, World!")
		}),
		tg.NewGroupCommand("mycounter").ActionFunc(func(c *tg.GC) {
			d := c.Session().Value.(*UserData)
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
		WithValue(&BotData{
			Name: "Jay",
		}).
		Debug(true)

	log.Printf("Authorized on account %s", bot.Api.Self.UserName)
	bot.Run()
}
