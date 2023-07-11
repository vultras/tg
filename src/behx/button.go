package behx

import (
    apix "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// The type wraps Telegram API's button to provide Action functionality.
type Button struct {
	Text string
	Data string
	Url string
	Action Action
}

type ButtonMap map[string] *Button

// Represents the reply button row.
type ButtonRow []*Button

// Returns new button with specified text and action.
func NewButton(text string, action Action) *Button {
	return &Button{
		Text: text,
		Action: action,
	}
}

func NewButtonData(text string, data string, action Action) *Button {
	return &Button{
		Text: text,
		Data: data,
		Action: action,
	}
}

func NewButtonUrl(text string, url string, action Action) *Button {
	return &Button{
		Text: text,
		Url: url,
		Action: action,
	}
}

func (btn *Button) ToTelegram() apix.KeyboardButton {
	return apix.NewKeyboardButton(btn.Text)
}

func (btn *Button) ToTelegramInline() apix.InlineKeyboardButton {
	if btn.Data != "" {
		return apix.NewInlineKeyboardButtonData(btn.Text, btn.Data)
	}
	
	if btn.Url != "" {
		return apix.NewInlineKeyboardButtonURL(btn.Text, btn.Url)
	}
	
	// If no match then return the data one with data the same as the text.
	return apix.NewInlineKeyboardButtonData(btn.Text, btn.Text)
}


func NewButtonRow(btns ...*Button) ButtonRow {
	return btns
}

