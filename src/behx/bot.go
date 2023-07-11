package behx

import (
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
    	BotAPI: bot,
		Behaviour: beh,
		sessions: make(SessionMap),
		
    }, nil
}

// Run the bot with the Behaviour.
func (bot *Bot) Run() error {
	bot.Debug = true
	
	uc := apix.NewUpdate(0)
	uc.Timeout = 60
	
	updates := bot.GetUpdatesChan(uc)
	
	for u := range updates {
		
		// Create new session if the one does not exist
		// for this user.
		sid := SessionId(u.Message.Chat.ID)
		if _, ok := bot.sessions[sid] ; !ok {
			bot.sessions.Add(sid)
		}
		
		session := bot.sessions[sid]
		ctx := &Context{
			B: bot,
			S: session,
		}
		
		// The "start" command resets the bot
		// by executing the Start Action.
		if u.Message.IsCommand() {
			cmd := u.Message.Command()
			if cmd == "start" {
				bot.Start.Act(ctx)
			}
			continue
		}
		
		screen := bot.Screens[session.CurrentScreenId]
		screen.Render(ctx)
	}
	
	return nil
}

