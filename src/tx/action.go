package tx

// Implementing the intereface lets you
// provide behaviour for the buttons etc.
type Action interface {
	Act(*Context)
}

// Customized action for the bot.
type ActionFunc func(*Context)

// The type implements changing screen to the underlying ScreenId
type ScreenChange ScreenId

func (sc ScreenChange) Act(c *Context) {
	if !c.B.ScreenExist(ScreenId(sc)) {
		panic(ScreenNotExistErr)
	}
	err := c.ChangeScreen(ScreenId(sc))
	if err != nil {
		panic(err)
	}
}

func (af ActionFunc) Act(c *Context) {
	af(c)
}
