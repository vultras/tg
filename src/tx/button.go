package tx

import (
	apix "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// The type wraps Telegram API's button to provide Action functionality.
type Button struct {
	Text         string
	Data         string
	Url          string
	SendLocation bool
	Action       *action
}

type ButtonMap map[string]*Button

// Represents the reply button row.
type ButtonRow []*Button

// Returns new button with the specified text and no action.
func NewButton(text string) *Button {
	return &Button{
		Text: text,
	}
}

// Set the URL for the button. Only for inline buttons.
func (btn *Button) WithUrl(url string) *Button {
	btn.Url = url
	return btn
}

// Set the action when pressing the button.
// By default is nil and does nothing.
func (btn *Button) WithAction(a Action) *Button {
	btn.Action = newAction(a)
	return btn
}

func (btn *Button) WithData(dat string) *Button {
	btn.Data = dat
	return btn
}

// Sets whether the button must send owner's location.
func (btn *Button) WithSendLocation(ok bool) *Button {
	btn.SendLocation = ok
	return btn
}

func (btn *Button) ActionFunc(fn ActionFunc) *Button {
	return btn.WithAction(fn)
}

func (btn *Button) ScreenChange(sc ScreenChange) *Button {
	return btn.WithAction(sc)
}

func (btn *Button) ToTelegram() apix.KeyboardButton {
	ret := apix.NewKeyboardButton(btn.Text)
	if btn.SendLocation {
		ret.RequestLocation = true
	}
	return ret
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

// Return the key of the button to identify it by messages and callbacks.
func (btn *Button) Key() string {
	if btn.Data != "" {
		return btn.Data
	}

	// If no match then return the data one with data the same as the text.
	return btn.Text
}

func NewButtonRow(btns ...*Button) ButtonRow {
	return btns
}
