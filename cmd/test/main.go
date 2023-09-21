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

func (w *MutateMessageWidget) Serve(c *tg.Context) {
	args, ok := c.Arg.(tg.ArgSlice)
	if ok {
		for _, arg := range args {
			c.Sendf("%v", arg)
		}
	}
	for u := range c.Input() {
		text := u.Message.Text
		c.Sendf("%s", w.Mutate(text))
	}
}

func (w *MutateMessageWidget) Filter(u *tg.Update, _ tg.MessageMap) bool {
	if u.Message == nil {
		return true
	}
	return false
}

func ExtractSessionData(c *tg.Context) *SessionData {
	return c.Session.Data.(*SessionData)
}

var (
	startScreenButton = tg.NewButton("Home").Go("/")
	backButton = tg.NewButton("Back").Go("..")
	backKeyboard = tg.NewKeyboard().Row(
		backButton,
	)

	incDecKeyboard = tg.NewKeyboard().Row(
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

	navKeyboard = tg.NewKeyboard().Row(
		tg.NewButton("Inc/Dec").Go("/inc-dec"),
	).Row(
		tg.NewButton("Mutate messages").Go("/mutate-messages"),
	).Row(
		tg.NewButton("Send location").Go("/send-location"),
	).Reply().WithOneTime(true)

	sendLocationKeyboard = tg.NewKeyboard().Row(
		tg.NewButton("Send location").
			WithSendLocation(true).
			ActionFunc(func(c *tg.Context) {
				l := c.Message.Location
				c.Sendf(
					"Longitude: %f\n"+
					"Latitude: %f\n"+
					"Heading: %d"+
					"",
					l.Longitude,
					l.Latitude,
					l.Heading,
				)
			}),
	).Row(
		startScreenButton,
	).Reply()

	// The keyboard to return to the start screen.
	navToStartKeyboard = tg.NewKeyboard().Row(
		startScreenButton,
	).Reply()
)

var beh = tg.NewBehaviour().
WithInitFunc(func(c *tg.Context) {
	// The session initialization.
	c.Session.Data = &SessionData{}
}).WithRootNode(tg.NewRootNode(
	// The "/" widget.
	tg.NewPage().
		WithInline(
			tg.NewKeyboard().Row(
				tg.NewButton("GoT Github page").
					WithUrl("https://github.com/mojosa-software/got"),
			).Inline().Widget("The bot started!"),
		).WithReply(
			navKeyboard.Widget("Choose what you are interested in"),
		),

	tg.NewNode(
		"mutate-messages", tg.NewPage().WithReply(
			tg.NewKeyboard().Row(
				tg.NewButton("Upper case").Go("upper-case"),
				tg.NewButton("Lower case").Go("lower-case"),
			).Row(
				backButton,
			).Reply().Widget(
				"Choose the function to mutate string",
			),
		),
		tg.NewNode(
			"upper-case", tg.NewPage().WithReply(
				backKeyboard.Reply().Widget(
					"Type a string and the bot will convert it to upper case",
				),
			).WithSub(
				NewMutateMessageWidget(strings.ToUpper),
			),
		),
		tg.NewNode(
			"lower-case", tg.NewPage().WithReply(
				backKeyboard.Reply().Widget(
					"Type a string and the bot will convert it to lower case",
				),
			).WithSub(
				NewMutateMessageWidget(strings.ToLower),
			),
		),
	),

	tg.NewNode(
		"inc-dec", tg.NewPage().WithReply(
				incDecKeyboard.Reply().Widget("Press the buttons to increment and decrement"),
			).ActionFunc(func(c *tg.Context) {
				// The function will be calleb before serving page.
				d := ExtractSessionData(c)
				c.Sendf("Current counter value = %d", d.Counter)
			}),
	),

	tg.NewNode(
		"send-location", tg.NewPage().WithReply(
			sendLocationKeyboard.Widget("Press the button to send your location!"),
		).WithInline(
			tg.NewKeyboard().Row(
				tg.NewButton(
					"Check",
				).WithData(
					"check",
				).ActionFunc(func(c *tg.Context) {
					d := ExtractSessionData(c)
					c.Sendf("Counter = %d", d.Counter)
				}),
			).Inline().Widget("Press the button to display your counter"),
		),
	),
)).WithCommands(
	tg.NewCommand("start").
		Desc("start or restart the bot or move to the start screen").
		ActionFunc(func(c *tg.Context){
			c.Go("/")
		}),
	tg.NewCommand("hello").
		Desc("sends the 'Hello, World!' message back").
		ActionFunc(func(c *tg.Context) {
			c.Sendf("Hello, World!")
		}),
	tg.NewCommand("read").
		Desc("reads a string and sends it back").
		WidgetFunc(func(c *tg.Context) {
			c.Sendf("Type text and I will send it back to you")
			for u := range c.Input() {
				if u.Message == nil {
					continue
				}
				c.Sendf("You typed %q", u.Message.Text)
				break
			}
			c.Sendf("Done")
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
	log.Println(beh.Screens)
	token := os.Getenv("BOT_TOKEN")

	bot, err := tg.NewBot(token)
	if err != nil {
		log.Panic(err)
	}
	bot = bot.
		WithBehaviour(beh).
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
