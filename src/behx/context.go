package behx

import (
	apix "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Update = apix.Update

// The type represents way to interact with user in
// handling functions. Is provided to Act() function always.
type Context struct {
	S *Session
	B *Bot
	U *Update
}

// Changes screen of user to the Id one.
func (c *Context) ChangeScreen(screenId ScreenId) error {
	if !c.B.ScreenExists(screenId) {
		return ScreenNotExistErr
	}
	
	c.S.PreviousScreenId = c.S.CurrentScreenId
	c.S.CurrentScreenId = screenId
	
	return nil
}

// Sends to the user specified text.
func (c *Context) Send(text string) error {
	
	return nil
}

