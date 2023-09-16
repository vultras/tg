package tg

import (
	"fmt"

	//tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Customized actions for the group behaviour.
type GroupAction interface {
	Act(*GroupContext)
}

// The handler function type.
type GroupActionFunc func(*GroupContext)

func (af GroupActionFunc) Act(a *GroupContext) {
	af(a)
}

type GC = GroupContext

func (c *GroupContext) Session() *GroupSession {
	session, _ := c.Bot.GetGroupSession(
		SessionId(c.SentFrom().ID),
	)
	return session
}

type GroupContext struct {
	*groupContext
	*Update
}

// Context for interaction inside groups.
type groupContext struct {
	Session *GroupSession
	Bot     *Bot
	updates chan *Update
}

func (c *groupContext) run(a GroupAction, u *Update) {
	go a.Act(&GroupContext{
		groupContext: c,
		Update:       u,
	})
}

func (c *groupContext) handleUpdateChan(updates chan *Update) {
	var act GroupAction
	beh := c.Bot.groupBehaviour
	for u := range updates {
		if u.Message != nil {
			msg := u.Message
			if msg.IsCommand() {
				cmdName := CommandName(msg.Command())

				// Skipping the commands sent not to us.
				withAt := msg.CommandWithAt()
				if len(cmdName) == len(withAt) {
					continue
				}

				atName := withAt[len(cmdName)+1:]
				if c.Bot.Me.UserName != atName {
					continue
				}
				cmd, ok := beh.Commands[cmdName]
				if !ok {
					// Some lack of command handling
					continue
				}
				act = cmd.Action
			}
		}
		if act != nil {
			c.run(act, u)
		}
	}
}

func (c *groupContext) Sendf(
	format string,
	v ...any,
) (*Message, error) {
	return c.Send(NewMessage(
		fmt.Sprintf(format, v...),
	))
}

// Sends into the chat specified values converted to strings.
func (c *groupContext) Send(v Sendable) (*Message, error) {
	return c.Bot.Send(c.Session.Id, v)
}
