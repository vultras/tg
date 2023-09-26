package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// The type represents keyboard to be emdedded into the messages (inline in Telegram terms).
type Inline struct {
	*Keyboard
}

// Convert the inline keyboard to markup for the tgbotapi.
func (kbd *Inline) ToApi() tgbotapi.InlineKeyboardMarkup {
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

// Transform the keyboard to widget with the specified text.
func (kbd *Inline) Compo(text string) *InlineCompo {
	ret := &InlinCompo{}
	ret.Inline = kbd
	ret.Text = text
	return ret
}

// The type implements message with an inline keyboard.
type InlineCompo struct {
	Compo
	Text string
	*Inline
}

// Implementing the Sendable interface.
func (widget *InlineCompo) SendConfig(
	c *Context,
) (*SendConfig) {
	var text string
	if widget.Text != "" {
		text = widget.Text
	} else {
		text = ">"
	}

	sid := c.Session.Id.ToApi()
	msgConfig := tgbotapi.NewMessage(sid, text)
	msgConfig.ReplyMarkup = widget.ToApi()

	ret := &SendConfig{}
	ret.Message = &msgConfig
	return ret
}

// Implementing the Server interface.
func (widget *InlineCompo) Serve(c *Context) {
	for u := range c.Input() {
		var act Action
		if u.CallbackQuery == nil {
			continue
		}
		cb := tgbotapi.NewCallback(
			u.CallbackQuery.ID,
			u.CallbackQuery.Data,
		)
		data := u.CallbackQuery.Data

		_, err := c.Bot.Api.Request(cb)
		if err != nil {
			//return err
			continue
		}

		btns := widget.ButtonMap()
		btn, ok := btns[data]
		if !ok {
			continue
		}
		if btn != nil {
			act = btn.Action
		} else if widget.Action != nil {
			act = widget.Action
		}
		c.Run(act, u)
	}
}

// Implementing the Filterer interface.
func (compo *InlineCompo) Filter(u *Update) bool {
	if widget == nil || u.CallbackQuery == nil {
		return true
	}

	if u.CallbackQuery.Message.MessageID != 
			compo.Message.MessageID {
		return true
	}

	return false
}

