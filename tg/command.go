package tg

import (
	//"flag"

	apix "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Message = apix.Message
type CommandName string

type Command struct {
	Name        CommandName
	Description string
	Action      *action
}
type CommandMap map[CommandName]*Command

func NewCommand(name CommandName) *Command {
	return &Command{
		Name: name,
	}
}

func (c *Command) WithAction(a Action) *Command {
	c.Action = newAction(a)
	return c
}

func (c *Command) ActionFunc(af ActionFunc) *Command {
	return c.WithAction(af)
}

func (c *Command) Desc(desc string) *Command {
	c.Description = desc
	return c
}

type GroupCommand struct {
	Name        CommandName
	Description string
	Action      GroupAction
}
type GroupCommandMap map[CommandName]*GroupCommand

func NewGroupCommand(name CommandName) *GroupCommand {
	return &GroupCommand{
		Name: name,
	}
}

func (cmd *GroupCommand) WithAction(a GroupAction) *GroupCommand {
	cmd.Action = a
	return cmd
}

func (cmd *GroupCommand) ActionFunc(fn GroupActionFunc) *GroupCommand {
	return cmd.WithAction(fn)
}

func (cmd *GroupCommand) Desc(desc string) *GroupCommand {
	cmd.Description = desc
	return cmd
}
