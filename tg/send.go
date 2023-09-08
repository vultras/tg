package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageId int64
type Image any

// Implementing the interface lets the
// value to be sent.
type Sendable interface {
	SendConfig(SessionId, *Bot) (*SendConfig, error)
}

type Renderable interface {
	Render(SessionId, *Bot) ([]*SendConfig, error)
}

// The type is used as an endpoint to send messages
// via bot.Send .
type SendConfig struct {
	// Simple message with text.
	// to add text use lower image
	// or see the ParseMode for tgbotapi .
	Message *tgbotapi.MessageConfig

	// The image to be sent.
	Image *tgbotapi.PhotoConfig
}

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

