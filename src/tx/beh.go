package tx

// The package implements
// behaviour for the Telegram bots.

// The type describes behaviour for the bot.
type Behaviour struct {
	Start     Action
	Screens   ScreenMap
	Keyboards KeyboardMap
}

// Returns new empty behaviour.
func NewBehaviour() *Behaviour {
	return &Behaviour{
		Screens:   make(ScreenMap),
		Keyboards: make(KeyboardMap),
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
