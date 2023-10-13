package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type InvoiceCompo struct {
	*MessageCompo
	tgbotapi.InvoiceConfig
}

func (compo *InvoiceCompo) SendConfig(
	sid SessionId, bot *Bot,
) (*SendConfig) {
	return nil
}

