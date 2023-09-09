package tg

// Unique identifier for the screen.
type ScreenId string

// Screen statement of the bot.
// Mostly what buttons to show.
type Screen struct {
	Id ScreenId
	// The widget to run when reaching the screen.
	Widget Widget

}

// Map structure for the screens.
type ScreenMap map[ScreenId]*Screen

// Returns the new screen with specified name and widget.
func NewScreen(id ScreenId, widget Widget) *Screen {
	return &Screen{
		Id: id,
		Widget: widget,
	}
}

