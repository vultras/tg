package tx

import (
	"reflect"
)

// Jsonable Action.
type action struct {
	Type   string
	Action Action
}

func newAction(a Action) *action {
	typ, ok := actionMapByReflect[reflect.TypeOf(a)]
	if !ok {
		panic(ActionNotDefinedErr)
	}

	return &action{
		Type:   typ,
		Action: a,
	}
}

func (a *action) Act(arg *A) {
	if a.Action != nil {
		a.Action.Act(arg)
	}
}

// Customized actions for the bot.
type Action interface {
	Act(*Arg)
}

// Customized actions for the
type GroupAction interface {
	Act(*GroupArg)
}

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
	if !c.B.behaviour.ScreenExist(screenId) {
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
	screen := c.B.behaviour.Screens[screenId]
	c.prevScreen = c.curScreen
	c.curScreen = screen
	screen.Render(c.Context)
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
