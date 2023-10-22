package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// The type represents reply keyboards.
type Reply struct {
	*Keyboard
	// If true will be removed after one press.
	OneTime bool
	// If true will remove the keyboard on send.
	Remove bool
}

// Set if we should remove current keyboard on the user side
// when sending the keyboard.
func (kbd *Reply) WithRemove(remove bool) *Reply {
	kbd.Remove = remove
	return kbd
}

// Set if the keyboard should be hidden after
// one of buttons is pressede.
func (kbd *Reply) WithOneTime(oneTime bool) *Reply{
	kbd.OneTime = oneTime
	return kbd
}

// Convert the Keyboard to the Telegram API type of reply keyboard.
func (kbd *Reply) ToApi() any {
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

// The type implements reply keyboard widget.
type ReplyCompo struct {
	*MessageCompo
	*Reply
}

// Implementing the sendable interface.
func (compo *ReplyCompo) SendConfig(
	sid SessionId, bot *Bot,
) (*SendConfig) {
	sendConfig := compo.MessageCompo.SendConfig(sid, bot)
	sendConfig.Message.ReplyMarkup = compo.Reply.ToApi()
	return sendConfig
}

func (compo *ReplyCompo) Filter(
	u *Update,
) bool {
	if compo == nil || u.Message == nil {
		return true
	}

	_, ok := compo.ButtonMap()[u.Message.Text]
	if !ok {
		if u.Message.Location != nil {
			locBtn := compo.ButtonMap().LocationButton()
			if locBtn == nil {
				return true
			}
		} else {
			return true
		}
	}
	return false
}

// Implementing the UI interface.
func (compo *ReplyCompo) Serve(c *Context) {
	for u := range c.Input() {
		var btn *Button
		text := u.Message.Text
		btns := compo.ButtonMap()

		btn, ok := btns[text]
		if !ok {
			if u.Message.Location != nil {
				btn = btns.LocationButton()
			}
		}

		if btn != nil {
			c.WithUpdate(u).Run(btn.Action)
		}
	}
}

