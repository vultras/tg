package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MessageId int64

// Implementing the interface provides
// way to define what message will be
// sent to the side of a user.
type Sendable interface {
	SendConfig(SessionId, *Bot) (*SendConfig)
	SetMessage(*Message)
}

type Errorer interface {
	Err() error
}

// The type is used as an endpoint to send messages
// via bot.Send .
type SendConfig struct {
	// Message with text and keyboard.
	Message *tgbotapi.MessageConfig

	// The image to be sent.
	Photo *tgbotapi.PhotoConfig
	Document *tgbotapi.DocumentConfig
	Location *tgbotapi.LocationConfig
	Error error
}


type MessageMap map[string] *Message

// Convert to the bot.Api.Send format.
func (config *SendConfig) ToApi() tgbotapi.Chattable {
	switch {
	case config.Message != nil :
		return *(config.Message)
	case config.Photo != nil :
		return *(config.Photo)
	case config.Location != nil :
		return *(config.Location)
	case config.Document != nil :
		return *(config.Document)
	}
	return nil
}

