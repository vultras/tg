package tx

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
	Id ScreenId

	// Text to be sent to the user when changing to the screen.
	Text ScreenText

	// The keyboard to be sent in the message part.
	InlineKeyboardId KeyboardId

	// Keyboard to be displayed on the screen.
	KeyboardId KeyboardId

	// Action called on the reaching the screen.
	Action Action
}

// Map structure for the screens.
type ScreenMap map[ScreenId]*Screen

// Returns the new screen with specified Text and Keyboard.
func NewScreen(id ScreenId) *Screen {
	return &Screen{
		Id: id,
	}
}

// Returns the screen with specified text printing on appearing.
func (s *Screen) WithText(text ScreenText) *Screen {
	s.Text = text
	return s
}

func (s *Screen) IKeyboard(kbdId KeyboardId) *Screen {
	s.InlineKeyboardId = kbdId
	return s
}

func (s *Screen) Keyboard(kbdId KeyboardId) *Screen {
	s.KeyboardId = kbdId
	return s
}

func (s *Screen) WithAction(a Action) *Screen {
	s.Action = a
	return s
}

func (s *Screen) ActionFunc(a ActionFunc) *Screen {
	return s.WithAction(a)
}

// Rendering the screen text to string to be sent or printed.
func (st ScreenText) String() string {
	return string(st)
}

// Renders output of the screen only to the side of the user.
func (s *Screen) Render(c *Context) error {
	id := c.Id.ToTelegram()

	msg := apix.NewMessage(id, s.Text.String())

	if s.InlineKeyboardId != "" {
		kbd, ok := c.B.behaviour.Keyboards[s.InlineKeyboardId]
		if !ok {
			return KeyboardNotExistErr
		}
		msg.ReplyMarkup = kbd.ToTelegramInline()
	}

	_, err := c.B.Send(msg)
	if err != nil {
		return err
	}

	msg = apix.NewMessage(id, ">")
	// Checking if we need to resend the keyboard.
	if s.KeyboardId != c.KeyboardId {
		// Remove keyboard by default.
		var tkbd any
		tkbd = apix.NewRemoveKeyboard(true)

		// Replace keyboard with the new one.
		if s.KeyboardId != "" {
			kbd, ok := c.B.behaviour.Keyboards[s.KeyboardId]
			if !ok {
				return KeyboardNotExistErr
			}
			tkbd = kbd.ToTelegram()
		}

		msg.ReplyMarkup = tkbd
		if _, err := c.B.Send(msg); err != nil {
			return err
		}
	}

	return nil
}
