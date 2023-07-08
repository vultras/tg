package behx

// The type represents way to interact with user in
// handling functions. Is provided to Act() function always.
type Context struct {
	S *Session
	B *Bot
}

func (c *Context) ChangeScreen(screen ScreenId) error {
	if !bot.ScreenExists(screenId) {
		return ScreenNotExistErr
	}
	
	return nil
}

func (c *Context) Send(text string) error {
	
	return nil
}

