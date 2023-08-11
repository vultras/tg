package main

import (
	"log"
	"os"
	"strings"

	"github.com/mojosa-software/got/src/tx"
)

var startScreenButton = tx.NewButton().
	WithText("🏠 To the start screen").
	ScreenChange("start")

var beh = tx.NewBehaviour().

	// The function will be called every time
	// the bot is started.
	OnStartFunc(func(c *tx.Context) {
		c.V["counter"] = new(int)
		c.ChangeScreen("start")
	}).WithKeyboards(

	// Increment/decrement keyboard.
	tx.NewKeyboard("inc/dec").Row(
		tx.NewButton().WithText("+").ActionFunc(func(c *tx.Context) {
			counter := c.V["counter"].(*int)
			*counter++
			c.Sendf("%d", *counter)
		}),
		tx.NewButton().WithText("-").ActionFunc(func(c *tx.Context) {
			counter := c.V["counter"].(*int)
			*counter--
			c.Sendf("%d", *counter)
		}),
	).Row(
		startScreenButton,
	),

	// The navigational keyboard.
	tx.NewKeyboard("nav").Row(
		tx.NewButton().WithText("Inc/Dec").ScreenChange("inc/dec"),
	).Row(
		tx.NewButton().WithText("Upper case").ScreenChange("upper-case"),
		tx.NewButton().WithText("Lower case").ScreenChange("lower-case"),
	),

	tx.NewKeyboard("istart").Row(
		tx.NewButton().WithText("GoT Github page").
			WithUrl("https://github.com/mojosa-software/got"),
	),

	// The keyboard to return to the start screen.
	tx.NewKeyboard("nav-start").Row(
		startScreenButton,
	),
).WithScreens(
	tx.NewScreen("start").
		WithText(
			"The bot started!"+
				" The bot is supposed to provide basic"+
				" understand of how the API works, so just"+
				" horse around a bit to guess everything out"+
				" by yourself!",
		).Keyboard("nav").
		IKeyboard("istart"),

	tx.NewScreen("inc/dec").
		WithText(
			"The screen shows how"+
				"user separated data works"+
				"by saving the counter for each of users"+
				"separately.",
		).
		Keyboard("inc/dec").
		// The function will be called when reaching the screen.
		ActionFunc(func(c *tx.Context) {
			counter := c.V["counter"].(*int)
			c.Sendf("Current counter value = %d", *counter)
		}),

	tx.NewScreen("upper-case").
		WithText("Type text and the bot will send you the upper case version to you").
		Keyboard("nav-start").
		ActionFunc(mutateMessage(strings.ToUpper)),

	tx.NewScreen("lower-case").
		WithText("Type text and the bot will send you the lower case version").
		Keyboard("nav-start").
		ActionFunc(mutateMessage(strings.ToLower)),
)

func mutateMessage(fn func(string) string) tx.ActionFunc {
	return func(c *tx.Context) {
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

func main() {
	token := os.Getenv("BOT_TOKEN")

	bot, err := tx.NewBot(token, beh, nil)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)
	bot.Run()
}
