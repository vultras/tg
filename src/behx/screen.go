package behx

// Unique identifier for the screen.
type ScreenId string

// Should be replaced with something that can be
// dinamicaly rendered. (WIP)
type ScreenText string

// Screen statement of the bot.
// Mostly what buttons to show.
type Screen struct {
	// Text to be sent to the user when changing to the screen.
	Text ScreenText
	// Keyboard to be displayed on the screen.
	Keyboard *Keyboard
}

// Map structure for the screens.
type ScreenMap Map[ScreenId, *Screen]
map[ScreenId] *Screen

// Returns the new screen with specified Text and Keyboard.
func NewScreen(text ScreenText, kbd *Keyboard) *Screen {
	return &Screen {
		Text: text,
		Keyboard: kbd,
	}
}


