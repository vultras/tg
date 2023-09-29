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
	Action      Action
	Widget Widget
}
type CommandMap map[CommandName]*Command

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

func (c *Command) WithWidget(w Widget) *Command {
	c.Widget = w
	return c
}

func (c *Command) WidgetFunc(fn Func) *Command {
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

func (c *Command) Go(pth Path, args ...any) *Command {
	return c.WithAction(ScreenGo{
		Path: pth,
		Args: args,
	})
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
type CommandCompo struct {
	PreStart Action
	Commands  CommandMap
	Usage Action
}

// Returns new empty CommandCompo.
func NewCommandCompo(cmds ...*Command) *CommandCompo {
	ret := (&CommandCompo{}).WithCommands(cmds...)
	//ret.Commands = make(CommandMap)
	return ret
}

// Set the commands to handle.
func (w *CommandCompo) WithCommands(cmds ...*Command) *CommandCompo {
	if w.Commands == nil {
		w.Commands = make(CommandMap)
	}
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
func (w *CommandCompo) WithPreStart(a Action) *CommandCompo {
	w.PreStart = a
	return w
}

// Set the prestart action with function.
func (w *CommandCompo) WithPreStartFunc(fn ActionFunc) *CommandCompo {
	return w.WithPreStart(fn)
}

// Set the usage action.
func (w *CommandCompo) WithUsage(a Action) *CommandCompo {
	w.Usage = a
	return w
}

// Set the usage action with function.
func (w *CommandCompo) WithUsageFunc(fn ActionFunc) *CommandCompo {
	return w.WithUsage(fn)
}

func (widget *CommandCompo) Filter(
	u *Update,
) bool {
	if u.Message == nil || !u.Message.IsCommand() {
		return false
	}

	return false
}

// Implementing server.
func (compo *CommandCompo) Serve(c *Context) {
	commanders := make(map[CommandName] BotCommander)
	for k, v := range compo.Commands {
		commanders[k] = v
	}
	c.Bot.SetCommands(
		tgbotapi.NewBotCommandScopeAllPrivateChats(),
		commanders,
	)

	var cmdUpdates *UpdateChan
	for u := range c.Input() {
		if c.Path() == "" && u.Message != nil {
			// Skipping and executing the preinit action
			// while we have the empty screen.
			// E. g. the session did not start.
			if !(u.Message.IsCommand() && u.Message.Command() == "start") {
				c.WithUpdate(u).Run(compo.PreStart)
				continue
			}
		}

		if u.Message != nil && u.Message.IsCommand() {
			// Command handling.
			cmdName := CommandName(u.Message.Command())
			cmd, ok := compo.Commands[cmdName]
			if !ok {
				c.WithUpdate(u).Run(compo.Usage)
				continue
			}

			c.WithUpdate(u).Run(cmd.Action)
			if cmd.Widget != nil {
				cmdUpdates.Close()
				cmdUpdates = c.WithUpdate(u).RunWidget(cmd.Widget)
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
