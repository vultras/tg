package tx

import (
	apix "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

/*
var otherKeyboard = tgbotapi.NewReplyKeyboard(
    tgbotapi.NewKeyboardButtonRow(
        tgbotapi.NewKeyboardButton("a"),
        tgbotapi.NewKeyboardButton("b"),
        tgbotapi.NewKeyboardButton("c"),
    ),
    tgbotapi.NewKeyboardButtonRow(
        tgbotapi.NewKeyboardButton("d"),
        tgbotapi.NewKeyboardButton("e"),
        tgbotapi.NewKeyboardButton("f"),
    ),
)*/

type KeyboardId string

// The type represents reply keyboard which
// is supposed to be showed on a Screen.
type Keyboard struct {
	Id   KeyboardId
	Rows []ButtonRow
}

type KeyboardMap map[KeyboardId]*Keyboard

// Return the new reply keyboard with rows as specified.
func NewKeyboard(id KeyboardId) *Keyboard {
	return &Keyboard{
		Id: id,
	}
}

// Adds a new button row to the current keyboard.
func (kbd *Keyboard) Row(btns ...*Button) *Keyboard {
	kbd.Rows = append(kbd.Rows, btns)
	return kbd
}

// Convert the Keyboard to the Telegram API type.
func (kbd *Keyboard) toTelegram() apix.ReplyKeyboardMarkup {
	rows := [][]apix.KeyboardButton{}
	for _, row := range kbd.Rows {
		buttons := []apix.KeyboardButton{}
		for _, button := range row {
			buttons = append(buttons, button.ToTelegram())
		}
		rows = append(rows, buttons)
	}

	return apix.NewReplyKeyboard(rows...)
}

func (kbd *Keyboard) toTelegramInline() apix.InlineKeyboardMarkup {
	rows := [][]apix.InlineKeyboardButton{}
	for _, row := range kbd.Rows {
		buttons := []apix.InlineKeyboardButton{}
		for _, button := range row {
			buttons = append(buttons, button.ToTelegramInline())
		}
		rows = append(rows, buttons)
	}

	return apix.NewInlineKeyboardMarkup(rows...)
}

// Returns the map of buttons. Used to define the Action.
func (kbd *Keyboard) buttonMap() ButtonMap {
	ret := make(ButtonMap)
	for _, vi := range kbd.Rows {
		for _, vj := range vi {
			ret[vj.Key()] = vj
		}
	}

	return ret
}
