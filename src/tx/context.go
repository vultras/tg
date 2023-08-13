package tx

import (
	"fmt"

	apix "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// The type represents way to interact with user in
// handling functions. Is provided to Act() function always.
type Context struct {
	*Session
	B       *Bot
	updates chan *Update

	// Is true if currently reading the Update.
	readingUpdate bool
}

// Goroutie function to handle each user.
func (c *Context) handleUpdateChan(updates chan *Update) {
	var act Action
	bot := c.B
	session := c.Session
	beh := bot.behaviour
	c.run(beh.Start, nil)
	for u := range updates {
		screen := bot.behaviour.Screens[session.CurrentScreenId]
		// The part is added to implement custom update handling.
		if u.Message != nil {
			if u.Message.IsCommand() && !c.readingUpdate {
				cmdName := CommandName(u.Message.Command())
				cmd, ok := beh.Commands[cmdName]
				if ok {
					act = cmd.Action
				} else {
				}
			} else {
				kbd := beh.Keyboards[screen.KeyboardId]
				btns := kbd.buttonMap()
				text := u.Message.Text
				btn, ok := btns[text]
				if !ok {
					if u.Message.Location != nil {
						for _, b := range btns {
							if b.SendLocation {
								btn = b
								ok = true
							}
						}
					} else if c.readingUpdate {
						// Skipping the update sending it to
						// the reading goroutine.
						c.updates <- u
						continue
					}
				}

				if ok {
					act = btn.Action
				}
			}
		} else if u.CallbackQuery != nil {
			cb := apix.NewCallback(u.CallbackQuery.ID, u.CallbackQuery.Data)
			data := u.CallbackQuery.Data

			_, err := bot.Request(cb)
			if err != nil {
				panic(err)
			}
			kbd := beh.Keyboards[screen.InlineKeyboardId]
			btns := kbd.buttonMap()
			btn, ok := btns[data]
			if !ok && c.readingUpdate {
				c.updates <- u
				continue
			}
			act = btn.Action
		}
		if act != nil {
			c.run(act, u)
		}
	}
}

func (c *Context) run(a Action, u *Update) {
	go a.Act(&A{
		Context: c,
		U:       u,
	})
}

// Returns the next update ignoring current screen.
func (c *Context) ReadUpdate() (*Update, error) {
	c.readingUpdate = true
	u := <-c.updates
	c.readingUpdate = false
	if u == nil {
		return nil, NotAvailableErr
	}

	return u, nil
}

// Returns the next text message that the user sends.
func (c *Context) ReadTextMessage() (string, error) {
	u, err := c.ReadUpdate()
	if err != nil {
		return "", err
	}
	if u.Message == nil {
		return "", WrongUpdateType{}
	}

	return u.Message.Text, nil
}

// Sends to the user specified text.
func (c *Context) Send(v ...any) error {
	msg := apix.NewMessage(c.Id.ToTelegram(), fmt.Sprint(v...))
	_, err := c.B.Send(msg)
	return err
}

// Sends the formatted with fmt.Sprintf message to the user.
func (c *Context) Sendf(format string, v ...any) error {
	return c.Send(fmt.Sprintf(format, v...))
}

// Context for interaction inside groups.
type GroupContext struct {
	*GroupSession
	B       *Bot
	updates chan *Update
}

func (c *GroupContext) run(a GroupAction, u *Update) {
	go a.Act(&GA{
		GroupContext: c,
		Update:       u,
	})
}

func (c *GroupContext) handleUpdateChan(updates chan *Update) {
	var act GroupAction
	beh := c.B.groupBehaviour
	for u := range updates {
		if u.Message != nil {
			msg := u.Message
			if msg.IsCommand() {
				cmdName := CommandName(msg.Command())

				// Skipping the commands sent not to us.
				atName := msg.CommandWithAt()[len(cmdName)+1:]
				if c.B.Me.UserName != atName {
					continue
				}
				cmd, ok := beh.Commands[cmdName]
				if !ok {
					// Some lack of command handling
					continue
				}
				act = cmd.Action
			}
		}
		if act != nil {
			c.run(act, u)
		}
	}
}

func (c *GroupContext) Sendf(format string, v ...any) error {
	return c.Send(fmt.Sprintf(format, v...))
}

// Sends into the chat specified values converted to strings.
func (c *GroupContext) Send(v ...any) error {
	msg := apix.NewMessage(c.Id.ToTelegram(), fmt.Sprint(v...))
	_, err := c.B.Send(msg)
	return err
}
