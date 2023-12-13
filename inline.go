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
		if row == nil {
			continue
		}
		buttons := []tgbotapi.InlineKeyboardButton{}
		for _, button := range row {
			if button == nil {
				continue
			}
			buttons = append(buttons, button.ToTelegramInline())
		}
		rows = append(rows, buttons)
	}

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// The type implements message with an inline keyboard.
type InlineCompo struct {
	*MessageCompo
	*Inline
}

// Implementing the Sendable interface.
func (compo *InlineCompo) SendConfig(
	sid SessionId, bot *Bot,
) (*SendConfig) {
	sendConfig := compo.MessageCompo.SendConfig(sid, bot)
	sendConfig.Message.ReplyMarkup = compo.Inline.ToApi()

	return sendConfig
}

func (compo *InlineCompo) Update(c *Context) {
	var edit tgbotapi.Chattable
	markup := compo.Inline.ToApi()
	ln := len(markup.InlineKeyboard)
	if ln == 0 || compo.Inline.Rows == nil {
		edit = tgbotapi.NewEditMessageText(
			c.Session.Id.ToApi(),
			compo.Message.MessageID,
			compo.Text,
		)
	} else {
		edit = tgbotapi.NewEditMessageTextAndMarkup(
			c.Session.Id.ToApi(),
			compo.Message.MessageID,
			compo.Text,
			markup,
		)
	}
	msg, _ := c.Bot.Api.Send(edit)
	compo.Message = &msg
}

// Implementing the Filterer interface.
func (compo *InlineCompo) Filter(u *Update) bool {
	if compo == nil || u.CallbackQuery == nil {
		return true
	}

	if u.CallbackQuery.Message.MessageID != 
			compo.Message.MessageID {
		return true
	}

	return false
}

// Implementing the Server interface.
func (widget *InlineCompo) Serve(c *Context) {
	btns := widget.ButtonMap()
	for u := range c.Input() {
		var act Action
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

		btn, ok := btns[data]
		if !ok {
			continue
		}
		if btn != nil {
			act = btn.Action
		} else if widget.Action != nil {
			act = widget.Action
		}
		c.WithUpdate(u).Run(act)
	}
}

