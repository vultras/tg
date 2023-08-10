package main

import (
	"log"
	"os"

	"github.com/mojosa-software/got/src/tx"
)

var navKeyboard = tx.NewKeyboard("nav").Row(
	tx.NewButton().WithText("Inc/Dec").ScreenChange("inc/dec"),
).Row(
	tx.NewButton().WithText("Upper case").ScreenChange("upper-case"),
	tx.NewButton().WithText("Lower case").ScreenChange("lower-case"),
)

var incKeyboard = tx.NewKeyboard("inc/dec").Row(
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
)

var startScreen = tx.NewScreen("start").
	WithText("The bot started!").
	Keyboard("nav")

var incScreen = tx.NewScreen("inc/dec").
	WithText("The screen shows how user separated data works").
	IKeyboard("inc/dec").
	Keyboard("nav")

var beh = tx.NewBehaviour().
	OnStartFunc(func(c *tx.Context) {
		// The function will be called every time
		// the bot is started.
		c.V["counter"] = new(int)
		c.ChangeScreen("start")
	}).WithKeyboards(
	navKeyboard,
	incKeyboard,
).WithScreens(
	startScreen,
	incScreen,
)

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
