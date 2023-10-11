package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Location = tgbotapi.Location

type LocationCompo struct {
	*MessageCompo
	Location
}

func (compo *LocationCompo) SendConfig(
	sid SessionId, bot *Bot,
) (*SendConfig) {
	cid := sid.ToApi()
	loc := tgbotapi.NewLocation(
		cid,
		compo.Latitude,
		compo.Longitude,
	)
	ret := &SendConfig{
		Location: &loc,
	}

	return ret
}

