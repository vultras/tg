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

// Context for interaction inside groups.
type GroupContext struct {
	*GroupSession
	B *Bot
}

// Goroutie function to handle each user.
func (c *Context) handleUpdateChan(updates chan *Update) {
	bot := c.B
	session := c.Session
	c.run(bot.Start, nil)
	for u := range updates {
		screen := bot.Screens[session.CurrentScreenId]
		// The part is added to implement custom update handling.
		if u.Message != nil {
			var act Action
			if u.Message.IsCommand() && !c.readingUpdate {
				cmdName := CommandName(u.Message.Command())
				cmd, ok := bot.Behaviour.Commands[cmdName]
				if ok {
					act = cmd.Action
				} else {
				}
			} else {
				kbd := bot.Keyboards[screen.KeyboardId]
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

			if act != nil {
				c.run(act, u)
			}
		} else if u.CallbackQuery != nil {
			cb := apix.NewCallback(u.CallbackQuery.ID, u.CallbackQuery.Data)
			data := u.CallbackQuery.Data

			_, err := bot.Request(cb)
			if err != nil {
				panic(err)
			}
			kbd := bot.Keyboards[screen.InlineKeyboardId]
			btns := kbd.buttonMap()
			btn, ok := btns[data]
			if !ok && c.readingUpdate {
				c.updates <- u
				continue
			}
			c.run(btn.Action, u)
		}
	}
}

func (c *Context) run(a Action, u *Update) {
	go a.Act(&A{c, u})
}

// Changes screen of user to the Id one.
func (c *Arg) ChangeScreen(screenId ScreenId) error {
	// Return if it will not change anything.
	if c.CurrentScreenId == screenId {
		return nil
	}

	if !c.B.ScreenExist(screenId) {
		return ScreenNotExistErr
	}

	// Stop the reading by sending the nil.
	if c.readingUpdate {
		c.updates <- nil
	}

	screen := c.B.Screens[screenId]
	screen.Render(c.Context)

	c.Session.ChangeScreen(screenId)
	c.KeyboardId = screen.KeyboardId

	if screen.Action != nil {
		c.run(screen.Action, c.U)
	}

	return nil
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
