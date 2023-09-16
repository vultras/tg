package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageId int64
type Image any

// Implementing the interface lets the
// value to be sent.
type Sendable interface {
	SendConfig(SessionId, *Bot) *SendConfig
}

type Renderable interface {
	Render(SessionId, *Bot) ([]*SendConfig)
}

type Errorer interface {
	Err() error
}

// The type is used as an endpoint to send messages
// via bot.Send .
type SendConfig struct {
	// The name will be used to store
	// the message in the map.
	Name string
	// Message with text and keyboard.
	Message *tgbotapi.MessageConfig

	// The image to be sent.
	Image *tgbotapi.PhotoConfig
	Error error
}

func (cfg *SendConfig) WithName(name string) *SendConfig {
	cfg.Name = name
	return cfg
}

type MessageMap map[string] *Message

// Convert to the bot.Api.Send format.
func (config *SendConfig) ToApi() tgbotapi.Chattable {
	if config.Message != nil {
		return *config.Message
	}

	if config.Image != nil {
		return *config.Image
	}
	return nil
}

