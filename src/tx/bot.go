package tx

import (
	//"fmt"

	"errors"

	apix "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Update = apix.Update
type Chat = apix.Chat
type User = apix.User

// The wrapper around Telegram API.
type Bot struct {
	*apix.BotAPI
	Me *User
	// Private bot behaviour.
	behaviour *Behaviour
	// Group bot behaviour.
	groupBehaviour *GroupBehaviour
	// Bot behaviour in channels.
	channelBehaviour *ChannelBehaviour
	sessions         SessionMap
	groupSessions    GroupSessionMap
}

// Return the new bot with empty sessions and behaviour.
func NewBot(token string) (*Bot, error) {
	bot, err := apix.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		BotAPI: bot,
	}, nil
}

func (bot *Bot) GetSessionValueBySid(
	sid SessionId,
) (any, bool) {
	v, ok := bot.sessions[sid]
	return v.V, ok
}

func (bot *Bot) GetGroupSessionValue(
	sid SessionId,
) (any, bool) {
	v, ok := bot.groupSessions[sid]
	return v.V, ok
}

func (b *Bot) WithBehaviour(beh *Behaviour) *Bot {
	b.behaviour = beh
	b.sessions = make(SessionMap)
	return b
}

func (b *Bot) WithSessions(sessions SessionMap) *Bot {
	b.sessions = sessions
	return b
}

func (b *Bot) WithGroupBehaviour(beh *GroupBehaviour) *Bot {
	b.groupBehaviour = beh
	b.groupSessions = make(GroupSessionMap)
	return b
}

func (b *Bot) WithGroupSessions(sessions GroupSessionMap) *Bot {
	b.groupSessions = sessions
	return b
}

// Run the bot with the Behaviour.
func (bot *Bot) Run() error {
	if bot.behaviour == nil &&
		bot.groupBehaviour == nil {
		return errors.New("no behaviour defined")
	}
	bot.Debug = true
	uc := apix.NewUpdate(0)
	uc.Timeout = 60
	updates := bot.GetUpdatesChan(uc)
	handles := make(map[string]chan *Update)

	if bot.behaviour != nil {
		chn := make(chan *Update)
		handles["private"] = chn
		go bot.handlePrivate(chn)
	}

	if bot.groupBehaviour != nil {
		chn := make(chan *Update)
		handles["group"] = chn
		handles["supergroup"] = chn
		go bot.handleGroup(chn)
	}

	me, _ := bot.GetMe()
	bot.Me = &me
	for u := range updates {
		chn, ok := handles[u.FromChat().Type]
		if !ok {
			continue
		}

		chn <- &u
	}

	return nil
}

// The function handles updates supposed for the private
// chat with the bot.
func (bot *Bot) handlePrivate(updates chan *Update) {
	chans := make(map[SessionId]chan *Update)
	var sid SessionId
	for u := range updates {
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
}

func (bot *Bot) handleGroup(updates chan *Update) {
	var sid SessionId
	chans := make(map[SessionId]chan *Update)
	for u := range updates {
		sid = SessionId(u.FromChat().ID)
		// If no session add new.
		if _, ok := bot.groupSessions[sid]; !ok {
			bot.groupSessions.Add(sid)
			session := bot.groupSessions[sid]
			ctx := &GroupContext{
				B:            bot,
				GroupSession: session,
				updates:      make(chan *Update),
			}
			chn := make(chan *Update)
			chans[sid] = chn
			go ctx.handleUpdateChan(chn)
		}

		chn := chans[sid]
		chn <- u
	}
}
