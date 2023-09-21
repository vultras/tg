package tg

// The package implements
// behaviour for the Telegram bots.

// The type describes behaviour for the bot in personal chats.
type Behaviour struct {
	Root Widget
	Init      *action
	Screens   ScreenMap
}

// Returns new empty behaviour.
func NewBehaviour() *Behaviour {
	return &Behaviour{
		Screens: make(ScreenMap),
	}
}

// Set the root widget. Mostly the CommandWidget is used.
func (b *Behaviour) WithRoot(root Widget) *Behaviour {
	b.Root = root
	return b
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

func (b *Behaviour) WithRootNode(node *RootNode) *Behaviour {
	b.Screens = node.ScreenMap()
	return b
}

// The function sets screens.
/*func (b *Behaviour) WithScreens(
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
}*/

// The function sets as the standard root widget CommandWidget
// and its commands..
func (b *Behaviour) WithCommands(cmds ...*Command) *Behaviour {
	b.Root = NewCommandWidget().
		WithCommands(cmds...).
		WithPreStartFunc(func(c *Context){
			c.Sendf("Please, use the /start command to start")
		}).WithUsageFunc(func(c *Context){
			c.Sendf("No such command")
		})

	return b
}

// Check whether the screen exists in the behaviour.
func (beh *Behaviour) PathExist(pth Path) bool {
	_, ok := beh.Screens[pth]
	return ok
}

// Returns the screen by it's ID.
func (beh *Behaviour) GetScreen(pth Path) *Screen {
	pth = pth.Clean()
	if !beh.PathExist(pth) {
		panic(ScreenNotExistErr)
	}

	screen := beh.Screens[pth]
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
