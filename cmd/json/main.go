package main

import (
	"encoding/json"
	"fmt"
	"os"

	//"strings"

	_ "github.com/mojosa-software/goscript/packages"

	"github.com/mojosa-software/goscript/env"
	"github.com/mojosa-software/goscript/vm"
	"github.com/mojosa-software/got/src/tx"
)

type UserData struct {
	Counter int
}

type Code string

func (c Code) Act(a *tx.A) {
	var err error
	fmt.Println("In Act")
	e := env.NewEnv()
	e.Define("a", a)
	e.Define("NotAvailableErr", tx.NotAvailableErr)
	e.Define("panic", func(v any) { panic(v) })
	err = e.DefineType("UserData", UserData{})
	if err != nil {
		panic(err)
	}

	_, err = vm.Execute(e, nil, string(c))
	if err != nil {
		panic(err)
	}
}

var startScreenButton = tx.NewButton("🏠 To the start screen").
	WithAction(Code(`
		a.ChangeScreen("start")
	`))

var beh = tx.NewBehaviour().

	// The function will be called every time
	// the bot is started.
	WithStart(Code(`
		a.V = new(UserData)
		a.ChangeScreen("start")
	`)).WithKeyboards(

	// Increment/decrement keyboard.
	tx.NewKeyboard("inc/dec").Row(
		tx.NewButton("+").WithAction(Code(`
			d = a.V
			d.Counter++
			a.Sendf("%d", d.Counter)
		`)),
		tx.NewButton("-").WithAction(Code(`
			d = a.V
			d.Counter--
			a.Sendf("%d", d.Counter)
		`)),
	).Row(
		startScreenButton,
	),

	// The navigational keyboard.
	tx.NewKeyboard("nav").Row(
		tx.NewButton("Inc/Dec").WithAction(Code(`a.ChangeScreen("inc/dec")`)),
	).Row(
		tx.NewButton("Upper case").WithAction(Code(`a.ChangeScreen("upper-case")`)),
		tx.NewButton("Lower case").WithAction(Code(`a.ChangeScreen("lower-case")`)),
	).Row(
		tx.NewButton("Send location").
			WithSendLocation(true).
			WithAction(Code(`
				err = nil
				if a.U.Message.Location != nil {
					l = a.U.Message.Location
					err = a.Sendf("Longitude: %f\nLatitude: %f\nHeading: %d", l.Longitude, l.Latitude, l.Heading)
				} else {
					err = a.Send("Somehow wrong location was sent")
				}
				if err != nil {
					a.Send(err)
				}
			`)),
	),

	tx.NewKeyboard("istart").Row(
		tx.NewButton("My Telegram").
			WithUrl("https://t.me/surdeus"),
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
			"The screen shows how "+
				"user separated data works "+
				"by saving the counter for each of users "+
				"separately. ",
		).
		Keyboard("inc/dec").
		// The function will be called when reaching the screen.
		WithAction(Code(`
			d = a.V
			a.Sendf("Current counter value = %d", d.Counter)
		`)),

	tx.NewScreen("upper-case").
		WithText("Type text and the bot will send you the upper case version to you").
		Keyboard("nav-start").
		WithAction(Code(`
			strings = import("strings")
			for {
				msg, err = a.ReadTextMessage()
				if err == NotAvailableErr {
					break
				} else if err != nil {
					panic(err)
				}
	
				err = a.Sendf("%s", strings.ToUpper(msg))
				if err != nil {
					panic(err)
				}
			}
		`)),

	tx.NewScreen("lower-case").
		WithText("Type text and the bot will send you the lower case version").
		Keyboard("nav-start").
		WithAction(Code(`
			strings = import("strings")
			for {
				msg, err = a.ReadTextMessage()
				if err == NotAvailableErr {
					break
				} else if err != nil {
					panic(err)
				}
	
				err = a.Sendf("%s", strings.ToLower(msg))
				if err != nil {
					panic(err)
				}
			}
		`)),
).WithCommands(
	tx.NewCommand("hello").
		Desc("sends the 'Hello, World!' message back").
		WithAction(Code(`
			a.Send("Hello, World!")
		`)),
	tx.NewCommand("read").
		Desc("reads a string and sends it back").
		WithAction(Code(`
			a.Send("Type some text:")
			msg, err = a.ReadTextMessage()
			if err != nil {
				return
			}
			a.Sendf("You typed %q", msg)
		`)),
)

func main() {
	bts, err := json.MarshalIndent(beh, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", bts)

	/*jBeh := &tx.Behaviour{}
	err = json.Unmarshal(bts, jBeh)
	if err != nil {
		panic(err)
	}*/

	bot, err := tx.NewBot(os.Getenv("BOT_TOKEN"), beh, nil)
	if err != nil {
		panic(err)
	}

	err = bot.Run()
	if err != nil {
		panic(err)
	}

}
