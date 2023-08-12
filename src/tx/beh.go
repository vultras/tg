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

// The type describes behaviour for the bot in group chats.
type GroupBehaviour struct {
	// Will be called on adding the bot to the group.
	//Add GroupAction
	// List of commands
	Commands CommandMap
}

// The type describes behaviour for the bot in channels.
type ChannelBehaviour struct {
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

func (b *Behaviour) OnStartChangeScreen(
	id ScreenId,
) *Behaviour {
	return b.WithStart(ScreenChange(id))
}

// The function sets keyboards.
func (b *Behaviour) WithKeyboards(
	kbds ...*Keyboard,
) *Behaviour {
	for _, kbd := range kbds {
		if kbd.Id == "" {
			panic("empty keyboard ID")
		}
		_, ok := b.Keyboards[kbd.Id]
		if ok {
			panic("duplicate keyboard IDs")
		}
		b.Keyboards[kbd.Id] = kbd
	}
	return b
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

// The function sets group commands.
/*func (b *Behaviour) WithGroupCommands(cmds ...*Command) *Behaviour {
	for _, cmd := range cmds {
		if cmd.Name == "" {
			panic("empty group command name")
		}
		_, ok := b.GroupCommands[cmd.Name]
		if ok {
			panic("duplicate group command definition")
		}
		b.GroupCommands[cmd.Name] = cmd
	}
	return b
}*/

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
