package tx

import (
	apix "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Unique identifier for the screen.
type ScreenId string

// Screen statement of the bot.
// Mostly what buttons to show.
type Screen struct {
	Id ScreenId
	// The text to be displayed when the screen is
	// reached.
	Text string
	// The keyboard to be sent in the message part.
	InlineKeyboard *Keyboard
	// Keyboard to be displayed on the screen.
	Keyboard *Keyboard
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
func (s *Screen) WithText(text string) *Screen {
	s.Text = text
	return s
}

func (s *Screen) WithInlineKeyboard(ikbd *Keyboard) *Screen {
	s.InlineKeyboard = ikbd
	return s
}

func (s *Screen) WithIKeyboard(ikbd *Keyboard) *Screen {
	return s.WithInlineKeyboard(ikbd)
}

func (s *Screen) WithKeyboard(kbd *Keyboard) *Screen {
	s.Keyboard = kbd
	return s
}

func (s *Screen) WithAction(a Action) *Screen {
	s.Action = a
	return s
}

func (s *Screen) ActionFunc(a ActionFunc) *Screen {
	return s.WithAction(a)
}

// Renders output of the screen only to the side of the user.
func (s *Screen) Render(c *Context) error {
	id := c.Id.ToTelegram()
	kbd := s.Keyboard
	iKbd := s.InlineKeyboard

	var ch [2]apix.Chattable
	var txt string

	// Screen text and inline keyboard.
	if s.Text != "" {
		txt = s.Text
	} else if iKbd != nil {
		if iKbd.Text != "" {
			txt = iKbd.Text
		} else {
			// Default to send the keyboard.
			txt = ">"
		}
	}
	if txt != "" {
		msg := apix.NewMessage(id, txt)
		if iKbd != nil {
			msg.ReplyMarkup = iKbd.toTelegramInline()
		} else if kbd != nil {
			msg.ReplyMarkup = kbd.toTelegramReply()
			if _, err := c.B.Send(msg); err != nil {
				return err
			}
			return nil
		} else {
			msg.ReplyMarkup = apix.NewRemoveKeyboard(true)
			if _, err := c.B.Send(msg); err != nil {
				return err
			}
			return nil
		}
		ch[0] = msg
	}

	// Screen text and reply keyboard.
	txt = ""
	if kbd != nil {
		if kbd.Text != "" {
			txt = kbd.Text
		} else {
			txt = ">"
		}
		msg := apix.NewMessage(id, txt)
		msg.ReplyMarkup = kbd.toTelegramReply()
		ch[1] = msg
	} else {
		// Removing keyboard if there is none.
		msg := apix.NewMessage(id, ">")
		msg.ReplyMarkup = apix.NewRemoveKeyboard(true)
		ch[1] = msg
	}

	for _, m := range ch {
		if m != nil {
			if _, err := c.B.Send(m); err != nil {
				return err
			}
		}
	}

	return nil
}
