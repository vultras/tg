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

// Transform the keyboard to widget with the specified text.
func (kbd *Reply) Widget(text string) *ReplyWidget {
	ret := &ReplyWidget{}
	ret.Reply = kbd
	ret.Text = text
	return ret
}

// The type implements reply keyboard widget.
type ReplyWidget struct {
	Text string
	*Reply
}

// Implementing the sendable interface.
func (widget *ReplyWidget) SendConfig(
	sid SessionId,
	bot *Bot,
) (*SendConfig) {
	if widget == nil {
		msgConfig := tgbotapi.NewMessage(sid.ToApi(), ">")
		msgConfig.ReplyMarkup = tgbotapi.NewRemoveKeyboard(true)
		return &SendConfig{
			Message: &msgConfig,
		}
	}
	var text string
	if widget.Text != "" {
		text = widget.Text
	} else {
		text = ">"
	}

	msgConfig := tgbotapi.NewMessage(sid.ToApi(), text)
	msgConfig.ReplyMarkup = widget.ToApi()

	ret := &SendConfig{}
	ret.Message = &msgConfig
	return ret
}

func (widget *ReplyWidget) Filter(
	u *Update,
	msgs MessageMap,
) bool {
	if widget == nil {
		return true
	}

	if u.Message == nil {
		return true
	}


	_, ok := widget.ButtonMap()[u.Message.Text]
	if !ok {
		if u.Message.Location != nil {
			locBtn := widget.ButtonMap().LocationButton()
			if locBtn == nil {
				return true
			}
		} else {
			return true
		}
	}
	return false
}

// Implementing the Widget interface.
func (widget *ReplyWidget) Serve(c *Context) {
	for u := range c.Input() {
		if u.Message == nil || u.Message.Text == "" {
			continue
		}
		var btn *Button
		text := u.Message.Text
		btns := widget.ButtonMap()

		btn, ok := btns[text]
		if !ok {
			if u.Message.Location != nil {
				btn = btns.LocationButton()
			}
		}

		if btn != nil {
			c.Run(btn.Action, u)
		}
	}
}

