package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// The general keyboard type used both in Reply and Inline.
type Keyboard struct {
	Rows []ButtonRow
}

type ReplyKeyboard struct {
	Keyboard
	// If true will be removed after one press.
	OneTime bool
	// If true will remove the keyboard on send.
	Remove bool
}

// The keyboard to be emdedded into the messages.
type InlineKeyboard struct {
	Keyboard
}

func NewInline() *InlineKeyboard {
	ret := &InlineKeyboard{}
	return ret
}

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
// Adds a new button row to the current keyboard.
func (kbd *ReplyKeyboard) Row(btns ...*Button) *ReplyKeyboard {
	// For empty row. We do not need that.
	if len(btns) < 1 {
		return kbd
	}
	kbd.Rows = append(kbd.Rows, btns)
	return kbd
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

func (kbd *ReplyKeyboard) WithRemove(remove bool) *ReplyKeyboard {
	kbd.Remove = remove
	return kbd
}

func (kbd *ReplyKeyboard) WithOneTime(oneTime bool) *ReplyKeyboard {
	kbd.OneTime = oneTime
	return kbd
}

// Returns the map of buttons. Used to define the Action.
func (kbd Keyboard) buttonMap() ButtonMap {
	ret := make(ButtonMap)
	for _, vi := range kbd.Rows {
		for _, vj := range vi {
			ret[vj.Key()] = vj
		}
	}

	return ret
}
