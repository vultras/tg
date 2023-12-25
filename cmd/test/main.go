package main

import (
	"log"
	"os"
	"strings"
	"fmt"

	"github.com/di4f/tg"
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
	args, ok := c.Arg().([]any)
	if ok {
		for _, arg := range args {
			c.Sendf("%v", arg)
		}
	}
	for u := range c.Input() {
		text := u.Message.Text
		_, err := c.Sendf2("%s", w.Mutate(text))
		if err != nil {
			c.Sendf("debug: %q", err)
		}
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
	homeButton   = tg.NewButton("Home").Go("/")
	backButton   = tg.NewButton("Back").Go("-")
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
		return tg.UI{
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
				).Row(
					tg.NewButton("Dynamic panel").Go("panel"),
				).Reply(),
			),

			tg.Func(func(c *tg.Context) {
				for u := range c.Input() {
					if u.EditedMessage != nil {
						c.Sendf2("The new message is `%s`", u.EditedMessage.Text)
					}
				}
			}),
		}
	}),

	tg.NewNode(
		"panel",
		tg.RenderFunc(func(c *tg.Context) tg.UI {
			var (
				n = 0
				ln = 4
				panel *tg.PanelCompo
			)
			
			panel = tg.NewMessage(
				"Some panel",
			).Panel(c, tg.RowserFunc(func(c *tg.Context) []tg.ButtonRow{
				btns := []tg.ButtonRow{
					tg.ButtonRow{tg.NewButton("Static shit")},
				}
				for i:=0 ; i<ln ; i++ {
					num := 1 + n * ln + i
					btns = append(btns, tg.ButtonRow{
						tg.NewButton("%d", num).WithAction(tg.Func(func(c *tg.Context){
							c.Sendf("%d", num*num)
						})),
						tg.NewButton("%d", num*num),
					})
				}
				btns = append(btns, tg.ButtonRow{
					tg.NewButton("Prev").WithAction(tg.ActionFunc(func(c *tg.Context){
						n--
						panel.Update(c)
					})),
					tg.NewButton("Next").WithAction(tg.ActionFunc(func(c *tg.Context){
						n++
						panel.Update(c)
					})),
				})

				return btns
			}))

			return tg.UI{
				panel,
				tg.NewMessage("").Reply(
					backKeyboard.Reply(),
				),
			}
		}),
	),

	tg.NewNode(
		"mutate-messages", tg.RenderFunc(func(c *tg.Context) tg.UI {
			return tg.UI{
				tg.NewMessage(
					"Choose the function to mutate string",
				).Reply(
					tg.NewKeyboard().Row(
						tg.NewButton("Upper case").Go("upper-case"),
						tg.NewButton("Lower case").Go("lower-case"),
						tg.NewButton("Escape chars").Go("escape"),
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
		tg.NewNode(
			"escape", tg.RenderFunc(func(c *tg.Context) tg.UI {
				return tg.UI{
					tg.NewMessage(
						"Type a string and the bot will escape characters in it",
					).Reply(
						backKeyboard.Reply(),
					),
					NewMutateMessageWidget(tg.Escape2),
				}
			}),
		),
	),

	tg.NewNode(
		"inc-dec", tg.RenderFunc(func(c *tg.Context) tg.UI {
			var (
				kbd *tg.InlineCompo
				//cntMsg *tg.MessageCompo
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
				tg.NewMessage("Use the reply keyboard to get back").Reply(
					backKeyboard.Reply(),
				),
			}
		}),
	),

	tg.NewNode(
		"send-location", tg.RenderFunc(func(c *tg.Context) tg.UI {
			return tg.UI{
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
	WithUsage(tg.Func(func(c *tg.Context) {
		c.Sendf("There is no such command %q", c.Message.Command())
	})).WithPreStart(tg.Func(func(c *tg.Context) {
	c.Sendf("Please, use /start ")
})).WithCommands(
	tg.NewCommand("info", "info desc").
		ActionFunc(func(c *tg.Context) {
			c.SendfHTML(`<a href="https://res.cloudinary.com/demo/image/upload/v1312461204/sample.jpg">cock</a><strong>cock</strong> die`)
		}),
	tg.NewCommand(
		"start",
		"start or restart the bot or move to the start screen",
	).Go("/"),
	tg.NewCommand("hello", "sends the 'Hello, World!' message back").
		ActionFunc(func(c *tg.Context) {
			c.Sendf("Hello, World!")
		}),
	tg.NewCommand("read", "reads a string and sends it back").
		WithWidget(
			tg.Func(func(c *tg.Context) {
				str := c.ReadString("Type a string and I will send it back")
				c.Sendf2("You typed `%s`", str)
			}),
		),
	tg.NewCommand("cat", "sends a sample image of cat").
		ActionFunc(func(c *tg.Context) {
			f, err := os.Open("media/cat.jpg")
			if err != nil {
				c.Sendf("err: %s", err)
				return
			}
			defer f.Close()
			photo := tg.NewFile(f).Photo().Name("cat.jpg").Caption("A cat!")
			c.Send(photo)
		}),
	tg.NewCommand("document", "sends a sample text document").
		ActionFunc(func(c *tg.Context) {
			f, err := os.Open("media/hello.txt")
			if err != nil {
				c.Sendf("err: %s", err)
				return
			}
			defer f.Close()
			doc := tg.NewFile(f).Document().Name("hello.txt").Caption("The document")
			c.Send(doc)
		}),
	tg.NewCommand("botname", "get the bot name").
		WithAction(tg.Func(func(c *tg.Context) {
			bd := c.Bot.Data.(*BotData)
			c.Sendf("My name is %q", bd.Name)
		})),
	tg.NewCommand("dynamic", "check of the dynamic work").
		WithWidget(tg.Func(func(c *tg.Context) {
		})),
	tg.NewCommand("history", "print go history").
		WithAction(tg.Func(func(c *tg.Context) {
			c.Sendf("%q", c.History())
		})),
	tg.NewCommand("washington", "send location of the Washington").
		WithAction(tg.Func(func(c *tg.Context) {
			c.Sendf("Washington location")
			c.Send(
				tg.NewMessage("").Location(
					47.751076, -120.740135,
				),
			)
		})),
	tg.NewCommand("invoice", "invoice check").
		WithAction(tg.Func(func(c *tg.Context) {
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
