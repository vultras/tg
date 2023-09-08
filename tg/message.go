package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageConfig struct {
	To SessionId
	ReplyTo MessageId
	Text string
	Inline *InlineKeyboard
	Reply *ReplyKeyboard
}

func NewMessage(to SessionId, text string) *MessageConfig {
	ret := &MessageConfig{}
	ret.To = to
	ret.Text = text
	return ret
}

func (config *MessageConfig) WithInline(
	inline *InlineKeyboard,
) *MessageConfig  {
	config.Inline = inline
	return config
}

func (config *MessageConfig) WithReply(
	reply *ReplyKeyboard,
) *MessageConfig {
	config.Reply = reply
	return config
}

func (config *MessageConfig) SendConfig(
	sid SessionId, bot *Bot,
) (*SendConfig, error) {
	var ret SendConfig
	msg := tgbotapi.NewMessage(config.To.ToApi(), config.Text)
	if config.Inline != nil {
		msg.ReplyMarkup = config.Inline.ToApi()
	}
	// Reply shades the inline.
	if config.Reply != nil {
		msg.ReplyMarkup = config.Reply.ToApi()
	}

	ret.Message = &msg
	return &ret, nil
}
