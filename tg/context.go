package tg

import (
	"fmt"

	apix "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Context for interaction inside groups.
type groupContext struct {
	*GroupSession
	*Bot
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
	beh := c.groupBehaviour
	for u := range updates {
		if u.Message != nil {
			msg := u.Message
			if msg.IsCommand() {
				cmdName := CommandName(msg.Command())

				// Skipping the commands sent not to us.
				atName := msg.CommandWithAt()[len(cmdName)+1:]
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

func (c *groupContext) Sendf(format string, v ...any) error {
	return c.Send(fmt.Sprintf(format, v...))
}

// Sends into the chat specified values converted to strings.
func (c *groupContext) Send(v ...any) error {
	msg := apix.NewMessage(c.Id.ToTelegram(), fmt.Sprint(v...))
	_, err := c.Bot.Send(msg)
	return err
}
