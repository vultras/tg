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
	rows []ButtonRow
}

// Return the new reply keyboard with rows as specified.
func NewKeyboard(rows ...ButtonRow) *Keyboard {
	return &Keyboard{
		rows: rows,
	}
}

// Convert the Keyboard to the Telegram API type.
func (kbd *Keyboard) ToTelegram() ReplyKeyboardMarkup {
	rows := [][]apix.KeyboardButton{}
	for _, row := range kbd.rows {
		buttons := []apix.KeyboardButton{}
		for _, button := range row {
			buttons = append(buttons, apix.NewKeyboardButton(button.text))
		}
		rows = append(rows, buttons)
	} 
	
	return apix.NewReplyKeyboard(rows)
}

func (kbd *Keyboard) ButtonMap() ButtonMap

