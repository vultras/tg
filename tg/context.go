package tg

import (
	"fmt"

	apix "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type context struct {
	*Session
	*Bot
	updates chan *Update
	// Is true if currently reading the Update.
	readingUpdate bool

	curScreen, prevScreen *Screen
}

// The type represents way to interact with user in
// handling functions. Is provided to Act() function always.

// Goroutie function to handle each user.
func (c *context) handleUpdateChan(updates chan *Update) {
	beh := c.behaviour

	if beh.Init != nil {
		c.run(beh.Init, nil)
	}
	for u := range updates {
		var act Action
		screen := c.curScreen
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
				kbd := screen.Keyboard
				if kbd == nil {
					if c.readingUpdate {
						c.updates <- u
					}
					continue
				}
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
			cb := apix.NewCallback(
				u.CallbackQuery.ID,
				u.CallbackQuery.Data,
			)
			data := u.CallbackQuery.Data

			_, err := c.Request(cb)
			if err != nil {
				panic(err)
			}
			kbd := screen.InlineKeyboard
			if kbd == nil {
				if c.readingUpdate {
					c.updates <- u
				}
				continue
			}

			btns := kbd.buttonMap()
			btn, ok := btns[data]
			if !ok && c.readingUpdate {
				c.updates <- u
				continue
			}
			if !ok {
				c.Sendf("%q", btns)
				continue
			}
			act = btn.Action
		}
		if act != nil {
			c.run(act, u)
		}
	}
}

func (c *context) run(a Action, u *Update) {
	go a.Act(&Context{
		context: c,
		Update:  u,
	})
}

// Returns the next update ignoring current screen.
func (c *context) ReadUpdate() (*Update, error) {
	c.readingUpdate = true
	u := <-c.updates
	c.readingUpdate = false
	if u == nil {
		return nil, NotAvailableErr
	}

	return u, nil
}

// Returns the next text message that the user sends.
func (c *context) ReadTextMessage() (string, error) {
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
func (c *context) Send(v ...any) error {
	msg := apix.NewMessage(c.Id.ToTelegram(), fmt.Sprint(v...))
	_, err := c.Bot.Send(msg)
	return err
}

// Sends the formatted with fmt.Sprintf message to the user.
func (c *context) Sendf(format string, v ...any) error {
	return c.Send(fmt.Sprintf(format, v...))
}

// Context for interaction inside groups.
type groupContext struct {
	*GroupSession
	*Bot
	updates chan *Update
}

func (c *groupContext) run(a GroupAction, u *Update) {
	go a.Act(&GroupContext{
		groupContext: c,
		Update:       u,
	})
}

func (c *groupContext) handleUpdateChan(updates chan *Update) {
	var act GroupAction
	beh := c.groupBehaviour
	for u := range updates {
		if u.Message != nil {
			msg := u.Message
			if msg.IsCommand() {
				cmdName := CommandName(msg.Command())

				// Skipping the commands sent not to us.
				atName := msg.CommandWithAt()[len(cmdName)+1:]
				if c.Bot.Me.UserName != atName {
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

func (c *groupContext) Sendf(format string, v ...any) error {
	return c.Send(fmt.Sprintf(format, v...))
}

// Sends into the chat specified values converted to strings.
func (c *groupContext) Send(v ...any) error {
	msg := apix.NewMessage(c.Id.ToTelegram(), fmt.Sprint(v...))
	_, err := c.Bot.Send(msg)
	return err
}
