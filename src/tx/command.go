package tx

import (
	//"flag"

	apix "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Message = apix.Message

type CommandName string

type CommandContext struct {
	// The field declares way to interact with the group chat in
	// general.
	*Context
	Message *Message
}

type CommandMap map[CommandName]*Command

type CommandHandlerFunc func(*CommandContext)
type CommandHandler interface {
	Run(*Context)
}

type Command struct {
	Name        CommandName
	Description string
	Action      Action
}

func NewCommand(name CommandName) *Command {
	return &Command{
		Name: name,
	}
}

func (c *Command) WithAction(a Action) *Command {
	c.Action = a
	return c
}

func (c *Command) ActionFunc(af ActionFunc) *Command {
	return c.WithAction(af)
}

func (c *Command) Desc(desc string) *Command {
	c.Description = desc
	return c
}
