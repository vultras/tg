package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Simple text message type.
type MessageConfig struct {
	ParseMode string
	Text string
}

// Return new message with the specified text.
func NewMessage(text string) *MessageConfig {
	ret := &MessageConfig{}
	ret.Text = text
	ret.ParseMode = tgbotapi.ModeMarkdown
	return ret
}

func (msg *MessageConfig) withParseMode(mode string) *MessageConfig{
	msg.ParseMode = mode
	return msg
}

// Set the default Markdown parsing mode.
func (msg *MessageConfig) MD() *MessageConfig {
	return msg.withParseMode(tgbotapi.ModeMarkdown)
}

func (msg *MessageConfig) MD2() *MessageConfig {
	return msg.withParseMode(tgbotapi.ModeMarkdownV2)
}


func (msg *MessageConfig) HTML() *MessageConfig {
	return msg.withParseMode(tgbotapi.ModeHTML)
}

func (config *MessageConfig) SendConfig(
	sid SessionId, bot *Bot,
) (*SendConfig) {
	var ret SendConfig
	msg := tgbotapi.NewMessage(sid.ToApi(), config.Text)
	ret.Message = &msg
	ret.Message.ParseMode = config.ParseMode
	return &ret
}
