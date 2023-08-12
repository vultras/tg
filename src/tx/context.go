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
	bot.Start.Act(c)
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
				// Sending wrong messages to
				// the currently reading goroutine.
				if !ok && c.readingUpdate {
					c.updates <- u
					continue
				}

				if ok {
					act = btn.Action
				}
			}

			if act != nil {
				c.run(act)
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
			btn := btns[data]
			btn.Action.Act(c)
		}
	}
}

func (c *Context) run(a Action) {
	go a.Act(c)
}

// Changes screen of user to the Id one.
func (c *Context) ChangeScreen(screenId ScreenId) error {
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
	screen.Render(c)

	c.Session.ChangeScreen(screenId)
	c.KeyboardId = screen.KeyboardId

	if screen.Action != nil {
		c.run(screen.Action)
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
