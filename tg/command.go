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
	Widget Widget
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

func (c *Command) WithWidget(w Widget) *Command {
	c.Widget = w
	return c
}

func (c *Command) WidgetFunc(fn WidgetFunc) *Command {
	return c.WithWidget(fn)
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

// The type is used to recognize commands and execute
// its actions and widgets .
type CommandWidget struct {
	PreStart Action
	Commands  CommandMap
	Usage Action
}

// Returns new empty CommandWidget.
func NewCommandWidget() *CommandWidget {
	ret := &CommandWidget{}
	ret.Commands = make(CommandMap)
	return ret
}

// Set the commands to handle.
func (w *CommandWidget) WithCommands(cmds ...*Command) *CommandWidget {
	for _, cmd := range cmds {
		if cmd.Name == "" {
			panic("empty command name")
		}
		_, ok := w.Commands[cmd.Name]
		if ok {
			panic("duplicate command definition")
		}
		w.Commands[cmd.Name] = cmd
	}
	return w
}

// Set the prestart action.
func (w *CommandWidget) WithPreStart(a Action) *CommandWidget {
	w.PreStart = a
	return w
}

// Set the prestart action with function.
func (w *CommandWidget) WithPreStartFunc(fn ActionFunc) *CommandWidget {
	return w.WithPreStart(fn)
}

// Set the usage action.
func (w *CommandWidget) WithUsage(a Action) *CommandWidget {
	w.Usage = a
	return w
}

// Set the usage action with function.
func (w *CommandWidget) WithUsageFunc(fn ActionFunc) *CommandWidget {
	return w.WithUsage(fn)
}

func (widget *Command) Filter(
	u *Update,
	msgs ...*Message,
) bool {
	/*if u.Message == nil || !u.Message.IsCommand() {
		return false
	}*/

	return false
}

func (widget *CommandWidget) Serve(c *Context) {
	commanders := make(map[CommandName] BotCommander)
	for k, v := range widget.Commands {
		commanders[k] = v
	}
	c.Bot.SetCommands(
		tgbotapi.NewBotCommandScopeAllPrivateChats(),
		commanders,
	)

	var cmdUpdates *UpdateChan
	for u := range c.Input() {
		if c.CurScreen() == "" && u.Message != nil {
			// Skipping and executing the preinit action
			// while we have the empty screen.
			// E. g. the session did not start.
			if !(u.Message.IsCommand() && u.Message.Command() == "start") {
				c.Run(widget.PreStart, u)
				continue
			}
		}

		if u.Message != nil && u.Message.IsCommand() {
			// Command handling.
			cmdName := CommandName(u.Message.Command())
			cmd, ok := widget.Commands[cmdName]
			if !ok {
				c.Run(widget.Usage, u)
				continue
			}

			c.Run(cmd.Action, u)
			if cmd.Widget != nil {
				cmdUpdates.Close()
				cmdUpdates = c.RunWidget(cmd.Widget)
			}
			continue
		}
		
		if !cmdUpdates.Closed()  {
			// Send to the commands channel if we are
			// executing one.
			cmdUpdates.Send(u)
		} else {
			c.Skip(u)
		}
	}
}
