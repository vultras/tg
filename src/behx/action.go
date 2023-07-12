package behx

// Implementing the intereface lets you
// provide behaviour for the buttons etc.
type Action interface {
	Act(*Context)
}

// Customized action for the bot.
type CustomAction func(*Context)

// The type implements changing screen to the underlying ScreenId
type ScreenChange ScreenId

// Returns new ScreenChange.
func NewScreenChange(screen string) ScreenChange {
	return ScreenChange(screen)
}

// Returns new CustomAction.
func NewCustomAction(fn func(*Context)) CustomAction {
	return CustomAction(fn)
}

func (sc ScreenChange) Act(c *Context) {
	err := c.ChangeScreen(ScreenId(sc))
	if err != nil {
		panic(err)
	}
}

func (ca CustomAction) Act(c *Context) {
	ca(c)
}

