package tg

// Interface to interact with the user.
type Context struct {
	*context
	*Update
}

// Customized actions for the bot.
type Action interface {
	Act(*Context)
}

type ActionFunc func(*Context)

func (af ActionFunc) Act(c *Context) {
	af(c)
}

// The type implements changing screen to the underlying ScreenId
type ScreenChange ScreenId

func (sc ScreenChange) Act(c *Context) {
	if !c.behaviour.ScreenExist(ScreenId(sc)) {
		panic(ScreenNotExistErr)
	}
	err := c.ChangeScreen(ScreenId(sc))
	if err != nil {
		panic(err)
	}
}

type C = Context

// Changes screen of user to the Id one.
func (c *Context) ChangeScreen(screenId ScreenId) error {
	if !c.behaviour.ScreenExist(screenId) {
		return ScreenNotExistErr
	}

	// Stop the reading by sending the nil,
	// since we change the screen and
	// current goroutine needs to be stopped.
	if c.readingUpdate {
		c.updates <- nil
	}

	// Getting the screen and changing to
	// then executing its action.
	screen := c.behaviour.Screens[screenId]
	c.prevScreen = c.curScreen
	c.curScreen = screen
	screen.Render(c.context)
	if screen.Action != nil {
		c.run(screen.Action, c.Update)
	}

	return nil
}

func (c *Context) SessionValue() any {
	v, _ := c.SessionValueBySid(c.Id)
	return v
}
