package tg

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

func (a *action) Act(c *Context) {
	if a.Action != nil {
		a.Action.Act(c)
	}
}

// The argument for handling in channenl behaviours.
type ChannelContext struct {
}
type CC = ChannelContext
type ChannelAction struct {
	Act (*ChannelContext)
}
