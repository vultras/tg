package main

import (
	"log"
	"os"
	"strings"
	"fmt"

	"github.com/mojosa-software/got/tg"
	//"math/rand"
	//"strconv"
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
	args, ok := c.Arg.([]any)
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

func (w *MutateMessageWidget) Filter(u *tg.Update) bool {
	if u.Message == nil {
		return true
	}
	return false
}

func ExtractSessionData(c *tg.Context) *SessionData {
	return c.Session.Data.(*SessionData)
}

var (
	homeButton = tg.NewButton("Home").Go("/")
	backButton = tg.NewButton("Back").Go("..")
	backKeyboard = tg.NewKeyboard().Row(
		backButton,
	)

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
		backButton,
	).Reply()
)

var beh = tg.NewBehaviour().
WithInitFunc(func(c *tg.Context) {
	// The session initialization.
	c.Session.Data = &SessionData{}
}).WithRootNode(tg.NewRootNode(
	// The "/" widget.
	tg.RenderFunc(func(c *tg.Context) tg.UI {
		return tg.UI {
			tg.NewMessage(fmt.Sprintf(
				fmt.Sprint(
					"Hello, %s!\n",
					"The testing bot started!\n",
					"You can see the basics of usage in the ",
					"cmd/test/main.go file!",
				),
				c.SentFrom().UserName,
			)).Inline(
				tg.NewKeyboard().Row(
					tg.NewButton("GoT Github page").
						WithUrl("https://github.com/mojosa-software/got"),
				).Inline(),
			),

			tg.NewMessage("Choose your interest").Reply(
				tg.NewKeyboard().Row(
						tg.NewButton("Inc/Dec").Go("/inc-dec"),
					).Row(
						tg.NewButton("Mutate messages").Go("/mutate-messages"),
					).Row(
						tg.NewButton("Send location").Go("/send-location"),
					).Reply(),
			),
		}
	}),

	tg.NewNode(
		"mutate-messages", tg.RenderFunc(func(c *tg.Context) tg.UI {
			return tg.UI{
				tg.NewMessage(
					"Choose the function to mutate string",
				).Reply(
					tg.NewKeyboard().Row(
						tg.NewButton("Upper case").Go("upper-case"),
						tg.NewButton("Lower case").Go("lower-case"),
					).Row(
						backButton,
					).Reply(),
				),
			}
		}),
		tg.NewNode(
			"upper-case", tg.RenderFunc(func(c *tg.Context) tg.UI {
				return tg.UI{
					tg.NewMessage(
						"Type a string and the bot will convert it to upper case",
					).Reply(
						backKeyboard.Reply(),
					),
					NewMutateMessageWidget(strings.ToUpper),
				}
			}),
		),
		tg.NewNode(
			"lower-case", tg.RenderFunc(func(c *tg.Context) tg.UI {
				return tg.UI{
					tg.NewMessage(
						"Type a string and the bot will convert it to lower case",
					).Reply(
						backKeyboard.Reply(),
					),
					NewMutateMessageWidget(strings.ToLower),
				}
			}),
		),
	),

	tg.NewNode(
		"inc-dec", tg.RenderFunc(func(c *tg.Context) tg.UI {
			var (
				kbd *tg.InlineCompo
				inline, std, onlyInc, onlyDec *tg.Inline
			)


			d := ExtractSessionData(c)
			format := "Press the buttons to increment and decrement.\n" +
				"Current counter value = %d"

			incBtn := tg.NewButton("+").ActionFunc(func(c *tg.Context) {
					d.Counter++
					kbd.Text = fmt.Sprintf(format, d.Counter)
					if d.Counter == 5 {
						kbd.Inline = onlyDec
					} else {
						kbd.Inline = std
					}
					kbd.Update(c)
				})
			decBtn := tg.NewButton("-").ActionFunc(func(c *tg.Context) {
					d.Counter--
					kbd.Text = fmt.Sprintf(format, d.Counter)
					if d.Counter == -5 {
						kbd.Inline = onlyInc
					} else {
						kbd.Inline = std
					}
					kbd.Update(c)
					//c.Sendf("%d", d.Counter)
				})

			onlyInc = tg.NewKeyboard().Row(incBtn).Inline()
			onlyDec = tg.NewKeyboard().Row(decBtn).Inline()
			std = tg.NewKeyboard().Row(incBtn, decBtn).Inline()

			if d.Counter == 5 {
				inline = onlyDec
			} else if d.Counter == -5 {
				inline = onlyInc
			} else {
				inline = std
			}

			kbd = tg.NewMessage(
				fmt.Sprintf(format, d.Counter),
			).Inline(inline)

			return tg.UI{
				kbd,
				tg.NewMessage("").Reply(
					backKeyboard.Reply(),
				),
			}
		}),
	),

	tg.NewNode(
		"send-location", tg.RenderFunc(func(c *tg.Context) tg.UI {
			return tg.UI {
				tg.NewMessage(
					"Press the button to display your counter",
				).Inline(
				tg.NewKeyboard().Row(
						tg.NewButton(
							"Check",
						).WithData(
							"check",
						).WithAction(tg.Func(func(c *tg.Context) {
							d := ExtractSessionData(c)
							c.Sendf("Counter = %d", d.Counter)
						})),
					).Inline(),
				),

				tg.NewMessage(
					"Press the button to send your location!",
				).Reply(
					sendLocationKeyboard,
				),
			}
		}),
	),
)).WithRoot(tg.NewCommandCompo().
WithPreStart(tg.Func(func(c *tg.Context){
	c.Sendf("Please, use /start ")
})).WithCommands(
	tg.NewCommand("info").
		ActionFunc(func(c *tg.Context){
			c.SendfHTML(`<a href="https://res.cloudinary.com/demo/image/upload/v1312461204/sample.jpg">cock</a><strong>cock</strong> die`)
		}),
	tg.NewCommand("start").
		Desc(
			"start or restart the bot or move to the start screen",
		).Go("/"),
	tg.NewCommand("hello").
		Desc("sends the 'Hello, World!' message back").
		ActionFunc(func(c *tg.Context) {
			c.Sendf("Hello, World!")
		}),
	tg.NewCommand("read").
		Desc("reads a string and sends it back").
		WithWidget(
			tg.Func(func(c *tg.Context){
				str := c.ReadString("Type a string and I will send it back")
				c.Sendf2("You typed `%s`", str)
			}),
		),
	tg.NewCommand("image").
		Desc("sends a sample image").
		ActionFunc(func(c *tg.Context) {
			img := tg.NewFile("media/cat.jpg").Image().Caption("A cat!")
			c.Send(img)
		}),
	tg.NewCommand("botname").
		Desc("get the bot name").
		WithAction(tg.Func(func(c *tg.Context) {
			bd := c.Bot.Data.(*BotData)
			c.Sendf("My name is %q", bd.Name)
		})),
	tg.NewCommand("dynamic").
		Desc("check of the dynamic work").
		WithWidget(tg.Func(func(c *tg.Context){
		})),
	))

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
