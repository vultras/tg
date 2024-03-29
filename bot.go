package tg

import (
	"errors"
	"sort"

	//"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)


type Chat = tgbotapi.Chat
type User = tgbotapi.User

// The wrapper around Telegram API.
type Bot struct {
	// Custom data value.
	Data any
	Api *tgbotapi.BotAPI
	Me  *User
	// Private bot behaviour.
	behaviour *Behaviour
	// Group bot behaviour.
	//groupBehaviour *GroupBehaviour
	// Bot behaviour in channels.
	//channelBehaviour *ChannelBehaviour
	contexts map[SessionId] *context
	sessions         SessionMap
	//groupSessions    GroupSessionMap
}

// Return the new bot with empty sessions and behaviour.
func NewBot(token string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	return &Bot{
		Api: bot,
		contexts: make(map[SessionId] *context),
	}, nil
}

func (bot *Bot) Debug(debug bool) *Bot {
	bot.Api.Debug = debug
	return bot
}

// Send the Renderable to the specified session client side.
// Can be used for both group and private sessions because
// SessionId represents both for chat IDs.
func (bot *Bot) Send(
	sid SessionId, v Sendable,
) (*Message, error) {
	config := v.SendConfig(sid, bot)
	if config.Error != nil {
		return nil, config.Error
	}

	msg, err := bot.Api.Send(config.ToApi())
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// Send to the session specified its ID raw chattable from the tgbotapi.
func (bot *Bot) SendRaw(
	sid SessionId, v tgbotapi.Chattable,
) (*Message, error) {
	msg, err := bot.Api.Send(v)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// Get session by its ID. Can be used for any scope
// including private, group and channel.
func (bot *Bot) GetSession(
	sid SessionId,
) (*Session, bool) {
	session, ok := bot.sessions[sid]
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

/*func (b *Bot) WithGroupBehaviour(beh *GroupBehaviour) *Bot {
	b.groupBehaviour = beh
	b.groupSessions = make(GroupSessionMap)
	return b
}

func (b *Bot) WithGroupSessions(sessions GroupSessionMap) *Bot {
	b.groupSessions = sessions
	return b
}*/

func (bot *Bot) DeleteCommands() {
	//tgbotapi.NewBotCommandScopeAllPrivateChats(),
	cfg := tgbotapi.NewDeleteMyCommands()
	bot.Api.Request(cfg)
}

// Setting the command on the user side.
func (bot *Bot) SetCommands(
	scope tgbotapi.BotCommandScope,
	cmdMap CommandMap,
) error {
	// First the private commands.
	names := []string{}
	for name := range cmdMap {
		names = append(names, string(name))
	}
	sort.Strings(names)

	cmds := []*Command{}
	for _, name := range names {
		cmds = append(
			cmds,
			cmdMap[CommandName(name)],
		)
	}

	botCmds := []tgbotapi.BotCommand{}
	for _, cmd := range cmds {
		botCmds = append(botCmds, cmd.ToApi())
	}

	//tgbotapi.NewBotCommandScopeAllPrivateChats(),
	cfg := tgbotapi.NewSetMyCommandsWithScope(
		scope,
		botCmds...,
	)

	_, err := bot.Api.Request(cfg)
	if err != nil {
		return err
	}
	return nil
}

// Run the bot with the Behaviour.
func (bot *Bot) Run() error {
	if bot.behaviour == nil {
		return errors.New("no behaviour defined")
	}

	if bot.behaviour != nil && bot.behaviour.Root == nil {
		return errors.New("the root widget is not set, cannot run")
	}

	uc := tgbotapi.NewUpdate(0)
	uc.Timeout = 10
	updates := bot.Api.GetUpdatesChan(uc)
	handles := make(map[string] chan *Update)

	if bot.behaviour != nil {
		chn := make(chan *Update)
		handles["private"] = chn
		go bot.handlePrivate(chn)
	}

	/*if bot.groupBehaviour != nil {
		commanders := make(map[CommandName] BotCommander)
		for k, v := range bot.groupBehaviour.Commands {
			commanders[k] = v
		}
		bot.SetCommands(
			tgbotapi.NewBotCommandScopeAllGroupChats(),
			commanders,
		)
		chn := make(chan *Update)
		handles["group"] = chn
		handles["supergroup"] = chn
		go bot.handleGroup(chn)
	}*/

	me, _ := bot.Api.GetMe()
	bot.Me = &me
	for up := range updates {
		u := &Update{
			Update: &up,
		}
		chn, ok := handles[u.FromChat().Type]
		if !ok {
			continue
		}

		chn <- u
	}

	return nil
}

// The function handles updates supposed for the private
// chat with the bot.
func (bot *Bot) handlePrivate(updates chan *Update) {
	var sid SessionId
	for u := range updates {
		sid = SessionId(u.FromChat().ID)
		ctx, ctxOk := bot.contexts[sid]
		if u.Message != nil && !ctxOk {

			session, sessionOk := bot.sessions[sid]
			if !sessionOk {
				// Creating session if we have none.
				session = bot.sessions.Add(sid, PrivateSessionScope)
			}
			session = bot.sessions[sid]

			// Create context on any message
			// if we have no one.
			ctx = &context{
				Bot:     bot,
				Session: session,
				updates: NewUpdateChan(),
			}
			 if !ctxOk {
				 bot.contexts[sid] = ctx
			 }

			go (&Context{
				context: ctx,
				Update: u,
				input: ctx.updates,
			}).serve()
			ctx.updates.Send(u)
			continue
		}

		if ctxOk {
			ctx.updates.Send(u)
		}
	}
}
/*
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
*/
