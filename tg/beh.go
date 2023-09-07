package tg

// The package implements
// behaviour for the Telegram bots.

// The type describes behaviour for the bot in personal chats.
type Behaviour struct {
	PreStart *action
	Init      *action
	Screens   ScreenMap
	Keyboards KeyboardMap
	Commands  CommandMap
}

// Returns new empty behaviour.
func NewBehaviour() *Behaviour {
	return &Behaviour{
		Screens:   make(ScreenMap),
		Keyboards: make(KeyboardMap),
		Commands:  make(CommandMap),
	}
}

// The Action will be called on session creation,
// not when starting or restarting the bot with the Start Action.
func (b *Behaviour) WithInit(a Action) *Behaviour {
	b.Init = newAction(a)
	return b
}

// Alias to WithInit to simplify behaviour definitions.
func (b *Behaviour) WithInitFunc(
	fn ActionFunc,
) *Behaviour {
	return b.WithInit(fn)
}

// Defines pre-start action.
// E. g. when the user has not type the "/start" command.
// Mostly used to send the "/start" command back
// with some warning.
func (b *Behaviour) WithPreStart(a Action) *Behaviour {
	b.PreStart = newAction(a)
	return b
}

// Alias for WithPreStart to be used with function inside.
func (b *Behaviour) WithPreStartFunc(fn ActionFunc) *Behaviour {
	return b.WithPreStart(fn)
}

// The function sets screens.
func (b *Behaviour) WithScreens(
	screens ...*Screen,
) *Behaviour {
	for _, screen := range screens {
		if screen.Id == "" {
			panic("empty screen ID")
		}
		_, ok := b.Screens[screen.Id]
		if ok {
			panic("duplicate keyboard IDs")
		}
		b.Screens[screen.Id] = screen
	}
	return b
}

// The function sets commands.
func (b *Behaviour) WithCommands(cmds ...*Command) *Behaviour {
	for _, cmd := range cmds {
		if cmd.Name == "" {
			panic("empty command name")
		}
		_, ok := b.Commands[cmd.Name]
		if ok {
			panic("duplicate command definition")
		}
		b.Commands[cmd.Name] = cmd
	}
	return b
}

// Check whether the screen exists in the behaviour.
func (beh *Behaviour) ScreenExist(id ScreenId) bool {
	_, ok := beh.Screens[id]
	return ok
}

// Returns the screen by it's ID.
func (beh *Behaviour) GetScreen(id ScreenId) *Screen {
	if !beh.ScreenExist(id) {
		panic(ScreenNotExistErr)
	}

	screen := beh.Screens[id]
	return screen
}

// The type describes behaviour for the bot in group chats.
type GroupBehaviour struct {
	Init GroupAction
	// List of commands
	Commands GroupCommandMap
}

// Returns new empty group behaviour object.
func NewGroupBehaviour() *GroupBehaviour {
	return &GroupBehaviour{
		Commands: make(GroupCommandMap),
	}
}

// Sets an Action for initialization on each group connected to the
// group bot.
func (b *GroupBehaviour) WithInitAction(a GroupAction) *GroupBehaviour {
	b.Init = a
	return b
}

// The method reciveies a function to be called on initialization of the
// bot group bot.
func (b *GroupBehaviour) InitFunc(fn GroupActionFunc) *GroupBehaviour {
	return b.WithInitAction(fn)
}

// The method sets group commands.
func (b *GroupBehaviour) WithCommands(
	cmds ...*GroupCommand,
) *GroupBehaviour {
	for _, cmd := range cmds {
		if cmd.Name == "" {
			panic("empty command name")
		}
		_, ok := b.Commands[cmd.Name]
		if ok {
			panic("duplicate command definition")
		}
		b.Commands[cmd.Name] = cmd
	}
	return b
}

// The type describes behaviour for the bot in channels.
type ChannelBehaviour struct {
}
