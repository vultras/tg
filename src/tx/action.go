package tx

//apix "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Action interface {
	Act(*Arg)
}

type GroupAction interface {
	Act(*GroupArg)
}

// Customized actions for the bot.

type ActionFunc func(*Arg)

func (af ActionFunc) Act(a *Arg) {
	af(a)
}

type GroupActionFunc func(*GroupArg)

func (af GroupActionFunc) Act(a *GroupArg) {
	af(a)
}

// The type implements changing screen to the underlying ScreenId
type ScreenChange ScreenId

func (sc ScreenChange) Act(c *Arg) {
	if !c.B.behaviour.ScreenExist(ScreenId(sc)) {
		panic(ScreenNotExistErr)
	}
	err := c.ChangeScreen(ScreenId(sc))
	if err != nil {
		panic(err)
	}
}

// The argument for handling.
type Arg struct {
	// Current context.
	*Context
	// The update that made the action to be called.
	U *Update
}
type A = Arg

// Changes screen of user to the Id one.
func (c *Arg) ChangeScreen(screenId ScreenId) error {
	// Return if it will not change anything.
	if c.CurrentScreenId == screenId {
		return nil
	}

	if !c.B.behaviour.ScreenExist(screenId) {
		return ScreenNotExistErr
	}

	// Stop the reading by sending the nil.
	if c.readingUpdate {
		c.updates <- nil
	}

	screen := c.B.behaviour.Screens[screenId]
	screen.Render(c.Context)

	c.Session.ChangeScreen(screenId)
	c.KeyboardId = screen.KeyboardId

	if screen.Action != nil {
		c.run(screen.Action, c.U)
	}

	return nil
}

// The argument for handling in group behaviour.
type GroupArg struct {
	*GroupContext
	*Update
}
type GA = GroupArg

func (a *GA) SentFromSid() SessionId {
	return SessionId(a.SentFrom().ID)
}

func (a *GA) GetSessionValue() any {
	v, _ := a.B.GetSessionValueBySid(a.SentFromSid())
	return v
}

// The argument for handling in channenl behaviours.
type ChannelArg struct {
}
type CA = ChannelArg
type ChannelAction struct {
	Act (*ChannelArg)
}

type JsonTyper interface {
	JsonType() string
}

type JsonAction struct {
	Type   string
	Action Action
}

func (ja JsonAction) UnmarshalJSON(bts []byte, ptr any) error {
	return nil
}
