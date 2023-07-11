package behx

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

// The type represents reply keyboard which
// is supposed to be showed on a Screen.
type Keyboard struct {
	Rows []ButtonRow
}

// Return the new reply keyboard with rows as specified.
func NewKeyboard(rows ...ButtonRow) *Keyboard {
	return &Keyboard{
		Rows: rows,
	}
}

// Convert the Keyboard to the Telegram API type.
func (kbd *Keyboard) ToTelegram() apix.ReplyKeyboardMarkup {
	rows := [][]apix.KeyboardButton{}
	for _, row := range kbd.Rows {
		buttons := []apix.KeyboardButton{}
		for _, button := range row {
			buttons = append(buttons, apix.NewKeyboardButton(button.Text))
		}
		rows = append(rows, buttons)
	} 
	
	return apix.NewReplyKeyboard(rows...)
}

// Returns the map of buttons. Used to define the Action.
func (kbd *Keyboard) buttonMap() ButtonMap {
	ret := make(ButtonMap)
	for _, vi := range kbd.Rows {
		for _, vj := range vi {
			ret[vj.Text] = vj
		}
	}
	
	return ret
}

