package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// The general keyboard type used both in Reply and Inline.
type Keyboard struct {
	// The action is called if there is no
	// defined action for the button.
	Action *action
	Rows []ButtonRow
	buttonMap ButtonMap
}

// The type represents reply keyboards.
type ReplyKeyboard struct {
	Keyboard
	// If true will be removed after one press.
	OneTime bool
	// If true will remove the keyboard on send.
	Remove bool
}

// The type represents keyboard to be emdedded into the messages.
type InlineKeyboard struct {
	Keyboard
}

// Returns new empty inline keyboard.
func NewInline() *InlineKeyboard {
	ret := &InlineKeyboard{}
	return ret
}

// Returns new empty reply keyboard.
func NewReply() *ReplyKeyboard {
	ret := &ReplyKeyboard {}
	return ret
}

// Adds a new button row to the current keyboard.
func (kbd *InlineKeyboard) Row(btns ...*Button) *InlineKeyboard {
	// For empty row. We do not need that.
	if len(btns) < 1 {
		return kbd
	}
	kbd.Rows = append(kbd.Rows, btns)
	return kbd
}

// Set default action for the buttons in keyboard.
func (kbd *InlineKeyboard) WithAction(a Action) *InlineKeyboard {
	kbd.Action = newAction(a)
	return kbd
}

// Alias to WithAction to simpler define actions.
func (kbd *InlineKeyboard) ActionFunc(fn ActionFunc) *InlineKeyboard {
	return kbd.WithAction(fn)
}

// Transform the keyboard to widget with the specified text.
func (kbd *InlineKeyboard) Widget(text string) *InlineKeyboardWidget {
	ret := &InlineKeyboardWidget{}
	ret.InlineKeyboard = kbd
	ret.Text = text
	return ret
}

// Adds a new button row to the current keyboard.
func (kbd *ReplyKeyboard) Row(btns ...*Button) *ReplyKeyboard {
	// For empty row. We do not need that.
	if len(btns) < 1 {
		return kbd
	}
	kbd.Rows = append(kbd.Rows, btns)
	return kbd
}

// Set default action for the keyboard.
func (kbd *ReplyKeyboard) WithAction(a Action) *ReplyKeyboard {
	kbd.Action = newAction(a)
	return kbd
}

// Alias to WithAction for simpler callback declarations.
func (kbd *ReplyKeyboard) ActionFunc(fn ActionFunc) *ReplyKeyboard {
	return kbd.WithAction(fn)
}

func (kbd *ReplyKeyboard) Widget(text string) *ReplyKeyboardWidget {
	ret := &ReplyKeyboardWidget{}
	ret.ReplyKeyboard = kbd
	ret.Text = text
	return ret
}

// Convert the Keyboard to the Telegram API type of reply keyboard.
func (kbd *ReplyKeyboard) ToApi() any {
	// Shades everything.
	if kbd.Remove {
		return tgbotapi.NewRemoveKeyboard(true)
	}

	rows := [][]tgbotapi.KeyboardButton{}
	for _, row := range kbd.Rows {
		buttons := []tgbotapi.KeyboardButton{}
		for _, button := range row {
			buttons = append(buttons, button.ToTelegram())
		}
		rows = append(rows, buttons)
	}

	if kbd.OneTime {
		return tgbotapi.NewOneTimeReplyKeyboard(rows...)
	}

	return tgbotapi.NewReplyKeyboard(rows...)
}

// Convert the inline keyboard to markup for the tgbotapi.
func (kbd *InlineKeyboard) ToApi() tgbotapi.InlineKeyboardMarkup {
	rows := [][]tgbotapi.InlineKeyboardButton{}
	for _, row := range kbd.Rows {
		buttons := []tgbotapi.InlineKeyboardButton{}
		for _, button := range row {
			buttons = append(buttons, button.ToTelegramInline())
		}
		rows = append(rows, buttons)
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// Set if we should remove current keyboard on the user side
// when sending the keyboard.
func (kbd *ReplyKeyboard) WithRemove(remove bool) *ReplyKeyboard {
	kbd.Remove = remove
	return kbd
}

// Set if the keyboard should be hidden after
// one of buttons is pressede.
func (kbd *ReplyKeyboard) WithOneTime(oneTime bool) *ReplyKeyboard {
	kbd.OneTime = oneTime
	return kbd
}

// Returns the map of buttons. Used to define the Action.
func (kbd Keyboard) ButtonMap() ButtonMap {
	if kbd.buttonMap != nil {
		return kbd.buttonMap
	}
	ret := make(ButtonMap)
	for _, vi := range kbd.Rows {
		for _, vj := range vi {
			ret[vj.Key()] = vj
		}
	}
	kbd.buttonMap = ret

	return ret
}

