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

type Code struct {
	Code string
	Add  int
}

func NewCode(code string) *Code {
	return &Code{
		Code: code,
	}
}

func (c *Code) Act(a *tx.A) {
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

	_, err = vm.Execute(e, nil, c.Code)
	if err != nil {
		panic(err)
	}
}

func main() {
	tx.DefineAction("goscript", &Code{})

	var startScreenButton = tx.NewButton("🏠 To the start screen").
		WithAction(NewCode(`
		a.ChangeScreen("start")
	`))

	var (
		incDecKeyboard = tx.NewKeyboard("").Row(
			tx.NewButton("+").WithAction(NewCode(`
			d = a.V
			d.Counter++
			a.Sendf("%d", d.Counter)
		`)),
			tx.NewButton("-").WithAction(NewCode(`
			d = a.V
			d.Counter--
			a.Sendf("%d", d.Counter)
		`)),
		).Row(
			startScreenButton,
		)

		// The navigational keyboard.
		navKeyboard = tx.NewKeyboard("").Row(
			tx.NewButton("Inc/Dec").WithAction(NewCode(`a.ChangeScreen("inc/dec")`)),
		).Row(
			tx.NewButton("Upper case").WithAction(NewCode(`a.ChangeScreen("upper-case")`)),
			tx.NewButton("Lower case").WithAction(NewCode(`a.ChangeScreen("lower-case")`)),
		).Row(
			tx.NewButton("Send location").
				WithSendLocation(true).
				WithAction(NewCode(`
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
		)

		inlineKeyboard = tx.NewKeyboard("").Row(
			tx.NewButton("My Telegram").
				WithUrl("https://t.me/surdeus"),
		)

		// The keyboard to return to the start screen.
		navToStartKeyboard = tx.NewKeyboard("nav-start").Row(
			startScreenButton,
		)
	)
	var beh = tx.NewBehaviour().
		// The function will be called every time
		// the bot is started.
		WithInit(NewCode(`
		a.V = new(UserData)
	`)).
		WithScreens(
			tx.NewScreen("start").
				WithText(
					"The bot started!"+
						" The bot is supposed to provide basic"+
						" understand of how the API works, so just"+
						" horse around a bit to guess everything out"+
						" by yourself!",
				).WithKeyboard(navKeyboard).
				WithIKeyboard(inlineKeyboard),

			tx.NewScreen("inc/dec").
				WithText(
					"The screen shows how "+
						"user separated data works "+
						"by saving the counter for each of users "+
						"separately. ",
				).
				WithKeyboard(incDecKeyboard).
				// The function will be called when reaching the screen.
				WithAction(NewCode(`
			d = a.V
			a.Sendf("Current counter value = %d", d.Counter)
		`)),

			tx.NewScreen("upper-case").
				WithText("Type text and the bot will send you the upper case version to you").
				WithKeyboard(navToStartKeyboard).
				WithAction(NewCode(`
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
				WithKeyboard(navToStartKeyboard).
				WithAction(NewCode(`
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
		tx.NewCommand("start").
			Desc("start or restart the bot").
			WithAction(NewCode(`
					a.ChangeScreen("start")
				`)),
		tx.NewCommand("hello").
			Desc("sends the 'Hello, World!' message back").
			WithAction(NewCode(`
				a.Send("Hello, World!")
			`)),
		tx.NewCommand("read").
			Desc("reads a string and sends it back").
			WithAction(NewCode(`
				a.Send("Type some text:")
				msg, err = a.ReadTextMessage()
				if err != nil {
					return
				}
				a.Sendf("You typed %q", msg)
			`)),
	)
	bts, err := json.MarshalIndent(beh, "", "\t")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s", bts)

	jBeh := &tx.Behaviour{}
	err = json.Unmarshal(bts, jBeh)
	if err != nil {
		panic(err)
	}

	bot, err := tx.NewBot(os.Getenv("BOT_TOKEN"))
	if err != nil {
		panic(err)
	}

	bot = bot.WithBehaviour(jBeh)

	err = bot.Run()
	if err != nil {
		panic(err)
	}

}
