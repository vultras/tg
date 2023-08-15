package tx

// The package implements
// behaviour for the Telegram bots.

// The type describes behaviour for the bot in personal chats.
type Behaviour struct {
	Start     Action
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

func (b *Behaviour) WithStart(a Action) *Behaviour {
	b.Start = a
	return b
}

func (b *Behaviour) OnStartFunc(
	fn ActionFunc,
) *Behaviour {
	return b.WithStart(fn)
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
