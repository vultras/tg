package tx

import (
	"fmt"

	apix "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	privateChans := make(map[SessionId]chan *Update)
	groupChans := make(map[SessionId]chan *Update)
	for u := range updates {
		var chatType string

		if u.Message != nil {
			chatType = u.Message.Chat.Type
		} else if u.CallbackQuery != nil {
			chatType = u.Message.Chat.Type
		}

		switch chatType {
		case "private":
			bot.handlePrivate(&u, privateChans)
		case "group", "supergroup":
			bot.handleGroup(&u, groupChans)
		}
	}

	return nil
}

// The function handles updates supposed for the private
// chat with the bot.
func (bot *Bot) handlePrivate(u *Update, chans map[SessionId]chan *Update) {
	var sid SessionId
	if u.Message != nil {
		msg := u.Message

		if bot.Debug {
			fmt.Printf("is command: %q\n", u.Message.IsCommand())
			fmt.Printf("command itself: %q\n", msg.Command())
			fmt.Printf("command arguments: %q\n", msg.CommandArguments())
			fmt.Printf("is to me: %q\n", bot.IsMessageToMe(*msg))
		}

		// Create new session if the one does not exist
		// for this user.
		sid = SessionId(u.Message.Chat.ID)
		if _, ok := bot.sessions[sid]; !ok {
			bot.sessions.Add(sid)
		}

		// The "start" command resets the bot
		// by executing the Start Action.
		if u.Message.IsCommand() {
			cmdName := CommandName(u.Message.Command())
			if cmdName == "start" {
				// Getting current session and context.
				session := bot.sessions[sid]
				ctx := &Context{
					B:       bot,
					Session: session,
					updates: make(chan *Update),
				}

				chn := make(chan *Update)
				chans[sid] = chn
				// Starting the goroutine for the user.
				go ctx.handleUpdateChan(chn)
			}
		}
	} else if u.CallbackQuery != nil {
		sid = SessionId(u.CallbackQuery.Message.Chat.ID)
	}
	chn, ok := chans[sid]
	// The bot MUST get the "start" command.
	// It will do nothing otherwise.
	if ok {
		chn <- u
	}
}

// Not implemented yet.
func (bot *Bot) handleGroup(u *Update, chans map[SessionId]chan *Update) {
}
