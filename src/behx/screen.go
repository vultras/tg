package behx

import (
	apix "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

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
	// The keyboard to be sent in the message part.
	InlineKeyboard *Keyboard
}

// Map structure for the screens.
type ScreenMap map[ScreenId] *Screen

// Returns the new screen with specified Text and Keyboard.
func NewScreen(text ScreenText, kbd *Keyboard) *Screen {
	return &Screen {
		Text: text,
		Keyboard: kbd,
	}
}

// Rendering the screen text to string to be sent or printed.
func (st ScreenText) String() string {
	return string(st)
}

// Renders output of the screen to the side of the user.
func (s *Screen) Render(c *Context) {
	id := c.S.Id.ToTelegram()
	msg := apix.NewMessage(id, s.Text.String())
	
	// First sending the inline keyboard.
	if s.InlineKeyboard != nil {
		msg.ReplyMarkup = s.InlineKeyboard.ToTelegram()
		if _, err := c.B.Send(msg) ; err != nil {
			panic(err)
		}
	}
	
	// Then sending the screen one.
	if s.Keyboard != nil {
		msg = apix.NewMessage(id, "check")
		msg.ReplyMarkup = s.Keyboard.ToTelegram()
		if _, err := c.B.Send(msg) ; err != nil {
			panic(err)
		}
	}
}

