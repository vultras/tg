package behx

import (
    apix "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// The type wraps Telegram API's button to provide Action functionality.
type Button struct {
	apix.KeyboardButton
	Action Action
}

type ButtonMap map[string] *Button

// Represents the reply button row.
type ButtonRow []*Button

// Returns new button with specified text and action.
func NewButton(text string, action Action) *Button {
	return &Button{
		KeyboardButton: apix.NewKeyboardButton(text),
		Action: action,
	}
}

func NewButtonRow(btns ...*Button) ButtonRow {
	return btns
}

