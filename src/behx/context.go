package behx

import (
	apix "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"fmt"
)

// The type represents way to interact with user in
// handling functions. Is provided to Act() function always.
type Context struct {
	*Session
	B *Bot
}

// Goroutie function to handle each user.
func (ctx *Context) handleUpdateChan(updates chan *Update) {
	bot := ctx.B
	session := ctx.Session
	bot.Start.Act(ctx)
	for u := range updates {
		screen := bot.Screens[session.CurrentScreenId]
		
		if u.Message != nil {
		
			kbd := bot.Keyboards[screen.KeyboardId]
			btns := kbd.buttonMap()
			text := u.Message.Text
			btn, ok := btns[text]
			
			// Skipping wrong text messages.
			if !ok {
				continue
			}
			
			btn.Action.Act(ctx)
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
			btn.Action.Act(ctx)
		}
	}
}

// Changes screen of user to the Id one.
func (c *Context) ChangeScreen(screenId ScreenId) error {
	// Return if it will not change anything.
	if c.CurrentScreenId == screenId {
		return nil
	}
	
	if !c.B.ScreenExists(screenId) {
		return ScreenNotExistErr
	}
	
	screen := c.B.Screens[screenId]
	screen.Render(c)
	
	c.Session.ChangeScreen(screenId)
	c.KeyboardId = screen.KeyboardId
	
	return nil
}

// Sends to the user specified text.
func (c *Context) Send(text string) error {
	msg := apix.NewMessage(c.Id.ToTelegram(), text)
	_, err := c.B.Send(msg)
	return err
}

func (c *Context) Sendf(format string, v ...any) error {
	return c.Send(fmt.Sprintf(format, v...))
}

