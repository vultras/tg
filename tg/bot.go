package tg

import (
	"errors"

	//"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Update = tgbotapi.Update
type Chat = tgbotapi.Chat
type User = tgbotapi.User

// The wrapper around Telegram API.
type Bot struct {
	Api *tgbotapi.BotAPI
	Me  *User
	// Private bot behaviour.
	behaviour *Behaviour
	// Group bot behaviour.
	groupBehaviour *GroupBehaviour
	// Bot behaviour in channels.
	channelBehaviour *ChannelBehaviour
	sessions         SessionMap
	groupSessions    GroupSessionMap
	value            any
	
}

// Return the new bot with empty sessions and behaviour.
func NewBot(token string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		Api: bot,
	}, nil
}

// Set the custom global value for the bot,
// so it can be accessed from the callback
// functions.
func (bot *Bot) WithValue(v any) *Bot {
	bot.value = v
	return bot
}

// Get the global bot value.
func (bot *Bot) Value() any {
	return bot.value
}

func (bot *Bot) Debug(debug bool) *Bot {
	bot.Api.Debug = debug
	return bot
}

func (bot *Bot) Send(
	sid SessionId, v Sendable,
) (*Message, error) {
	config, err := v.SendConfig(sid, bot)
	if err != nil {
		return nil, err
	}

	msg, err := bot.Api.Send(config.ToApi())
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func (bot *Bot) Render(
	sid SessionId, r Renderable,
) ([]*Message, error) {
	configs, err := r.Render(sid, bot)
	if err != nil {
		return []*Message{}, err
	}
	messages := []*Message{}
	for _, config := range configs {
		msg, err := bot.Api.Send(config.ToApi())
		if err != nil {
			return messages, err
		}
		messages = append(messages, &msg)
	}
	return messages, nil
}

func (bot *Bot) GetSession(
	sid SessionId,
) (*Session, bool) {
	session, ok := bot.sessions[sid]
	return session, ok
}

func (bot *Bot) GetGroupSession(
	sid SessionId,
) (*GroupSession, bool) {
	session, ok := bot.groupSessions[sid]
	return session, ok
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
	uc := tgbotapi.NewUpdate(0)
	uc.Timeout = 60
	updates := bot.Api.GetUpdatesChan(uc)
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

	me, _ := bot.Api.GetMe()
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
		sid = SessionId(u.FromChat().ID)
		// Create new session if the one does not exist
		// for this user.

		// Making the bot ignore anything except "start"
		// before the session started
		session, sessionOk := bot.sessions[sid]
		chn, chnOk := chans[sid]
		if sessionOk {
			// Creating new goroutine for 
			// the session that exists
			// but has none.
			if !chnOk {
				ctx := &context{
					Bot:     bot,
					Session: session,
				}
				chn := make(chan *Update)
				chans[sid] = chn
				go ctx.handleUpdateChan(chn)
			}
		} else if u.Message != nil {
			// Create session on any message
			// if we have no one.
			bot.sessions.Add(sid)
			lsession := bot.sessions[sid]
			ctx := &context{
				Bot:     bot,
				Session: lsession,
			}
			chn := make(chan *Update)
			chans[sid] = chn
			go ctx.handleUpdateChan(chn)
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
			ctx := &groupContext{
				Bot:     bot,
				Session: session,
				updates: make(chan *Update),
			}
			chn := make(chan *Update)
			chans[sid] = chn
			go ctx.handleUpdateChan(chn)
		}

		chn := chans[sid]
		chn <- u
	}
}
