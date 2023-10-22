package tg

// The package implements
// behaviour for the Telegram bots.

// The type describes behaviour for the bot in personal chats.
type Behaviour struct {
	Root Component
	Init      Action
	Screens   ScreenMap
}

// Returns new empty behaviour.
func NewBehaviour() *Behaviour {
	return &Behaviour{
		Screens: make(ScreenMap),
	}
}

// The Action will be called on session creation,
// not when starting or restarting the bot with the Start Action.
func (b *Behaviour) WithInit(a Action) *Behaviour {
	b.Init = a
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
func (b *Behaviour) WithRoot(root Component) *Behaviour {
	b.Root = root
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

