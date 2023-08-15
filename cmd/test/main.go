package main

import (
	"log"
	"os"
	"strings"

	"github.com/mojosa-software/got/src/tx"
)

type UserData struct {
	Counter int
}

var (
	startScreenButton = tx.NewButton("🏠 To the start screen").
				ScreenChange("start")

	incDecKeyboard = tx.NewKeyboard("").Row(
		tx.NewButton("+").ActionFunc(func(c *tx.A) {
			d := c.V.(*UserData)
			d.Counter++
			c.Sendf("%d", d.Counter)
		}),
		tx.NewButton("-").ActionFunc(func(c *tx.A) {
			d := c.V.(*UserData)
			d.Counter--
			c.Sendf("%d", d.Counter)
		}),
	).Row(
		startScreenButton,
	)

	navKeyboard = tx.NewKeyboard("Choose your interest").
			WithOneTime(true).
			Row(
			tx.NewButton("Inc/Dec").ScreenChange("inc/dec"),
		).Row(
		tx.NewButton("Upper case").ScreenChange("upper-case"),
		tx.NewButton("Lower case").ScreenChange("lower-case"),
	).Row(
		tx.NewButton("Send location").ScreenChange("send-location"),
	)

	sendLocationKeyboard = tx.NewKeyboard("Press the button to send your location").
				Row(
			tx.NewButton("Send location").
				WithSendLocation(true).
				ActionFunc(func(c *tx.A) {
					var err error
					if c.U.Message.Location != nil {
						l := c.U.Message.Location
						err = c.Sendf(
							"Longitude: %f\n"+
								"Latitude: %f\n"+
								"Heading: %d"+
								"",
							l.Longitude,
							l.Latitude,
							l.Heading,
						)
					} else {
						err = c.Send("Somehow wrong location was sent")
					}
					if err != nil {
						c.Send(err)
					}
				}),
		).Row(
		startScreenButton,
	)

	// The keyboard to return to the start screen.
	navToStartKeyboard = tx.NewKeyboard("").Row(
		startScreenButton,
	)
)

var beh = tx.NewBehaviour().
	OnStartFunc(func(c *tx.A) {
		// The function will be called every time
		// the bot is started.
		c.V = &UserData{}
		c.ChangeScreen("start")
	}).WithScreens(
	tx.NewScreen("start").
		WithText(
			"The bot started!"+
				" The bot is supposed to provide basic"+
				" understand of how the API works, so just"+
				" horse around a bit to guess everything out"+
				" by yourself!",
		).WithKeyboard(navKeyboard).
		// The inline keyboard with link to GitHub page.
		WithIKeyboard(
			tx.NewKeyboard("istart").Row(
				tx.NewButton("GoT Github page").
					WithUrl("https://github.com/mojosa-software/got"),
			),
		),

	tx.NewScreen("inc/dec").
		WithText(
			"The screen shows how "+
				"user separated data works "+
				"by saving the counter for each of users "+
				"separately. ",
		).
		WithKeyboard(incDecKeyboard).
		// The function will be called when reaching the screen.
		ActionFunc(func(c *tx.A) {
			d := c.V.(*UserData)
			c.Sendf("Current counter value = %d", d.Counter)
		}),

	tx.NewScreen("upper-case").
		WithText("Type text and the bot will send you the upper case version to you").
		WithKeyboard(navToStartKeyboard).
		ActionFunc(mutateMessage(strings.ToUpper)),

	tx.NewScreen("lower-case").
		WithText("Type text and the bot will send you the lower case version").
		WithKeyboard(navToStartKeyboard).
		ActionFunc(mutateMessage(strings.ToLower)),

	tx.NewScreen("send-location").
		WithText("Send your location and I will tell where you are!").
		WithKeyboard(sendLocationKeyboard).
		WithIKeyboard(
			tx.NewKeyboard("").Row(
				tx.NewButton("Check").
					WithData("check").
					ActionFunc(func(a *tx.A) {
						d := a.V.(*UserData)
						a.Sendf("Counter = %d", d.Counter)
					}),
			),
		),
).WithCommands(
	tx.NewCommand("hello").
		Desc("sends the 'Hello, World!' message back").
		ActionFunc(func(c *tx.A) {
			c.Send("Hello, World!")
		}),
	tx.NewCommand("read").
		Desc("reads a string and sends it back").
		ActionFunc(func(c *tx.A) {
			c.Send("Type some text:")
			msg, err := c.ReadTextMessage()
			if err != nil {
				return
			}
			c.Sendf("You typed %q", msg)
		}),
)

func mutateMessage(fn func(string) string) tx.ActionFunc {
	return func(c *tx.A) {
		for {
			msg, err := c.ReadTextMessage()
			if err == tx.NotAvailableErr {
				break
			} else if err != nil {
				panic(err)
			}

			err = c.Sendf("%s", fn(msg))
			if err != nil {
				panic(err)
			}
		}
	}
}

var gBeh = tx.NewGroupBehaviour().
	InitFunc(func(a *tx.GA) {
	}).
	WithCommands(
		tx.NewGroupCommand("hello").ActionFunc(func(a *tx.GA) {
			a.Send("Hello, World!")
		}),
		tx.NewGroupCommand("mycounter").ActionFunc(func(a *tx.GA) {
			d := a.GetSessionValue().(*UserData)
			a.Sendf("Your counter value is %d", d.Counter)
		}),
	)

func main() {
	token := os.Getenv("BOT_TOKEN")

	bot, err := tx.NewBot(token)
	if err != nil {
		log.Panic(err)
	}
	bot = bot.
		WithBehaviour(beh).
		WithGroupBehaviour(gBeh)

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	bot.Run()
}
