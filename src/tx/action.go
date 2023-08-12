package tx

import (
	apix "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Update = apix.Update

// The argument for handling.
type Arg struct {
	// Current context.
	*Context
	// The update that made the action to be called.
	U *Update
}
type A = Arg

type GroupArg struct {
	GroupArg *GroupContext
	U        *Update
}
type GA = GroupArg

type Action interface {
	Act(*Arg)
}

// Customized action for the bot.
type ActionFunc func(*Arg)

// The type implements changing screen to the underlying ScreenId
type ScreenChange ScreenId

func (sc ScreenChange) Act(c *Arg) {
	if !c.B.ScreenExist(ScreenId(sc)) {
		panic(ScreenNotExistErr)
	}
	err := c.ChangeScreen(ScreenId(sc))
	if err != nil {
		panic(err)
	}
}

func (af ActionFunc) Act(c *Arg) {
	af(c)
}
