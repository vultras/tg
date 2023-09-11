package tg

import (
	//"flag"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotCommander interface {
	ToApi() tgbotapi.BotCommand
}
type Message = tgbotapi.Message
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

func (c *Command) ToApi() tgbotapi.BotCommand {
	ret := tgbotapi.BotCommand{}
	ret.Command = string(c.Name)
	ret.Description = c.Description
	return ret
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

func (c *GroupCommand) ToApi() tgbotapi.BotCommand {
	ret := tgbotapi.BotCommand{}
	ret.Command = string(c.Name)
	ret.Description = c.Description
	return ret
}
