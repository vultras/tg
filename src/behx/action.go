package behx

// Implementing the intereface lets you
// provide behaviour for the buttons etc.
type Action interface {
	Act(*Bot)
}

// Customized action for the bot.
type CustomAction func(*Bot)

// The type implements changing screen to the underlying ScreenId
type ScreenChange ScreenId

// Returns new ScreenChange.
func NewScreenChange(screen string) ScreenChange {
	return ScreenChange(screen)
}

// Returns new CustomAction.
func NewCustomAction(fn func(*Bot)) CustomAction {
	return CustomAction(fn)
}

func (sc ScreenChange) Act(c *Context) {
	c.ChangeScreen(ScreenId(sc))
}

func (ca CustomAction) Act(c *Context) {
	ca(bot)
}


/*
// The type describes interface
// defining what should be done.
type Actioner interface {
	Act(*bot)
}

// Simple way to define that the button should just change
// screen to the Id.
type ScreenChanger ScreenId

// Custom function to be executed on button press.
type ActionFunc func(*Bot)

func (a ActionFunc) Act(bot *Bot) {
	a(bot)
}

func (sc ScreenChanger) Act(bot *Bot) {
	bot.ChangeScreenTo(sc)
}
*/

