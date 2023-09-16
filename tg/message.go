package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Simple text message type.
type MessageConfig struct {
	Text string
}

// Return new message with the specified text.
func NewMessage(text string) *MessageConfig {
	ret := &MessageConfig{}
	ret.Text = text
	return ret
}

func (config *MessageConfig) SendConfig(
	sid SessionId, bot *Bot,
) (*SendConfig) {
	var ret SendConfig
	msg := tgbotapi.NewMessage(sid.ToApi(), config.Text)
	ret.Message = &msg
	return &ret
}
