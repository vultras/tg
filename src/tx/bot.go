package tx

import (
	apix "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	//"log"
)

// The wrapper around Telegram API.
type Bot struct {
	*apix.BotAPI
	*Behaviour
	sessions SessionMap
}

// Return the new bot for running the Behaviour.
func NewBot(token string, beh *Behaviour, sessions SessionMap) (*Bot, error) {
	bot, err := apix.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	// Make new sessions if no current are provided.
	if sessions == nil {
		sessions = make(SessionMap)
	}

	return &Bot{
		BotAPI:    bot,
		Behaviour: beh,
		sessions:  make(SessionMap),
	}, nil
}

// Run the bot with the Behaviour.
func (bot *Bot) Run() error {
	bot.Debug = true

	uc := apix.NewUpdate(0)
	uc.Timeout = 60

	updates := bot.GetUpdatesChan(uc)

	chans := make(map[SessionId]chan *Update)
	for u := range updates {
		var sid SessionId
		if u.Message != nil {
			// Create new session if the one does not exist
			// for this user.
			sid = SessionId(u.Message.Chat.ID)
			if _, ok := bot.sessions[sid]; !ok {
				bot.sessions.Add(sid)
			}

			// The "start" command resets the bot
			// by executing the Start Action.
			if u.Message.IsCommand() {
				cmd := u.Message.Command()
				if cmd == "start" {
					// Getting current session and context.
					session := bot.sessions[sid]
					ctx := &Context{
						B:       bot,
						Session: session,
					}

					chn := make(chan *Update)
					chans[sid] = chn
					// Starting the goroutine for the user.
					go ctx.handleUpdateChan(chn)
					continue
				}
			}
		} else if u.CallbackQuery != nil {
			sid = SessionId(u.CallbackQuery.Message.Chat.ID)
		}
		chn, ok := chans[sid]
		if ok {
			chn <- &u
		}
	}

	return nil
}
