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
	Inline *InlineKeyboard
	// Keyboard to be displayed on the screen.
	Reply *ReplyKeyboard
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

func (s *Screen) WithInline(ikbd *InlineKeyboard) *Screen {
	s.Inline= ikbd
	return s
}

func (s *Screen) WithReply(kbd *ReplyKeyboard) *Screen {
	s.Reply = kbd
	return s
}

func (s *Screen) WithAction(a Action) *Screen {
	s.Action = newAction(a)
	return s
}

func (s *Screen) ActionFunc(a ActionFunc) *Screen {
	return s.WithAction(a)
}

func (s *Screen) Render(
	sid SessionId, bot *Bot,
) ([]*SendConfig, error) {
	cid := sid.ToApi()
	reply := s.Reply
	inline := s.Inline
	ret := []*SendConfig{}
	var txt string
	// Screen text and inline keyboard.
	if s.Text != "" {
		txt = s.Text
	} else if inline != nil {
		// Default to send the keyboard.
		txt = ">"
	}
	if txt != "" {
		msgConfig := tgbotapi.NewMessage(cid, txt)
		if inline != nil {
			msgConfig.ReplyMarkup = inline.ToApi()
		} else if reply != nil {
			msgConfig.ReplyMarkup = reply.ToApi()
			ret = append(ret, &SendConfig{Message: &msgConfig})
			return ret, nil
		} else {
			msgConfig.ReplyMarkup = NewReply().
				WithRemove(true).
				ToApi()
			ret = append(ret, &SendConfig{Message: &msgConfig})
			return ret, nil
		}
		ret = append(ret, &SendConfig{Message: &msgConfig})
	}

	// Screen text and reply keyboard.
	if reply != nil {
		msgConfig := tgbotapi.NewMessage(cid, ">")
		msgConfig.ReplyMarkup = reply.ToApi()
		ret = append(ret, &SendConfig{
			Message: &msgConfig,
		})
	} else {
		// Removing keyboard if there is none.
		msgConfig := tgbotapi.NewMessage(cid, ">")
		msgConfig.ReplyMarkup = NewReply().
			WithRemove(true).
			ToApi()
			ret = append(ret, &SendConfig{Message: &msgConfig})
	}

	return ret, nil
}
