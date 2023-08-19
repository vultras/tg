package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	Action *action
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
	s.Action = newAction(a)
	return s
}

func (s *Screen) ActionFunc(a ActionFunc) *Screen {
	return s.WithAction(a)
}

// Renders output of the screen only to the side of the user.
func (s *Screen) Render(
	sid SessionId, bot *Bot,
) ([]*Message, error) {
	cid := sid.ToApi()
	kbd := s.Keyboard
	iKbd := s.InlineKeyboard

	var ch [2]tgbotapi.Chattable
	var txt string

	msgs := []*Message{}

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
		msgConfig := tgbotapi.NewMessage(cid, txt)
		if iKbd != nil {
			msgConfig.ReplyMarkup = iKbd.toTelegramInline()
		} else if kbd != nil {
			msgConfig.ReplyMarkup = kbd.toTelegramReply()
			msg, err := bot.Api.Send(msgConfig)
			if err != nil {
				return msgs, err
			}
			msgs = append(msgs, &msg)
			return msgs, nil
		} else {
			msgConfig.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
			msg, err := bot.Api.Send(msgConfig)
			if err != nil {
				return msgs, err
			}
			msgs = append(msgs, &msg)
			return msgs, nil
		}
		ch[0] = msgConfig
	}

	// Screen text and reply keyboard.
	txt = ""
	if kbd != nil {
		if kbd.Text != "" {
			txt = kbd.Text
		} else {
			txt = ">"
		}
		msgConfig := tgbotapi.NewMessage(cid, txt)
		msgConfig.ReplyMarkup = kbd.toTelegramReply()
		ch[1] = msgConfig
	} else {
		// Removing keyboard if there is none.
		msgConfig := tgbotapi.NewMessage(cid, ">")
		msgConfig.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		ch[1] = msgConfig
	}

	for _, m := range ch {
		if m != nil {
			msg, err := bot.Api.Send(m)
			if err != nil {
				return msgs, err
			}
			msgs = append(msgs, &msg)
		}
	}

	return msgs, nil
}
