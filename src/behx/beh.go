package behx

// The package implements
// behaviour for the Telegram bots.


// The type describes behaviour for the bot.
type Behaviour struct {
	Start Action
	Screens ScreenMap
}


// Check whether the screen exists in the behaviour.
func (beh *Behaviour) ScreenExists(id ScreenId) bool {
	_, ok := bot.behaviour.Screens[id]
	return ok
}

// Returns the screen by it's ID.
func (beh *Behaviour) GetScreen(id ScreenId) *Screen {
	if !beh.ScreenExists(id) {
		panic(ScreenNotExistErr)
	}
	
	screen := beh.Screens[id]
	return screen
}

