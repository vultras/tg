package tg

import (
	"fmt"

	//tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type context struct {
	Session *Session
	// To reach the bot abilities inside callbacks.
	Bot     *Bot
	widgetUpdates chan *Update
	CurScreen, PrevScreen *Screen
}

// The type represents way to interact with user in
// handling functions. Is provided to Act() function always.

// Goroutie function to handle each user.
func (c *context) handleUpdateChan(updates chan *Update) {
	beh := c.Bot.behaviour

	session := c.Session
	preStart := beh.PreStart
	if beh.Init != nil {
		c.run(beh.Init, nil)
	}
	for u := range updates {
		// The part is added to implement custom update handling.
		if !session.started {
			if u.Message.IsCommand() &&
					u.Message.Command() == "start" {
				// Special treatment for the "/start"
				// command.
				session.started = true
				cmdName := CommandName("start")
				cmd, ok := beh.Commands[cmdName]
				if ok {
					if cmd.Action != nil {
						c.run(cmd.Action, u)
					}
				} else {
					// Some usage.
				}
			} else {
				// Prestart handling.
				c.run(preStart, u)
			}

			continue
		}

		if u.Message != nil && u.Message.IsCommand() {
			// Command handling.
			cmdName := CommandName(u.Message.Command())
			cmd, ok := beh.Commands[cmdName]
			if ok {
				if cmd.Action != nil {
					c.run(cmd.Action, u)
				}
			} else {
				// Some usage.
			}
			continue
		} 
		
		// The standard thing - send messages to widgets.
		c.widgetUpdates <- u
	}
}

func (c *context) run(a Action, u *Update) {
	go a.Act(&Context{
		context: c,
		Update:  u,
	})
}

func (c *context) Render(v Renderable) ([]*Message, error) {
	return c.Bot.Render(c.Session.Id, v)
}

// Sends to the Sendable object.
func (c *context) Send(v Sendable) (*Message, error) {
	return c.Bot.Send(c.Session.Id, v)
}

// Sends the formatted with fmt.Sprintf message to the user.
func (c *context) Sendf(format string, v ...any) (*Message, error) {
	msg, err := c.Send(NewMessage(
		c.Session.Id, fmt.Sprintf(format, v...),
	))
	if err != nil {
		return nil, err
	}
	return msg, err
}

// Interface to interact with the user.
type Context struct {
	*context
	// The update that called the Context usage.
	*Update
	// Used as way to provide outer values redirection
	// into widgets and actions 
	Arg any
}

// Customized actions for the bot.
type Action interface {
	Act(*Context)
}

type ActionFunc func(*Context)

func (af ActionFunc) Act(c *Context) {
	af(c)
}

// The type implements changing screen to the underlying ScreenId
type ScreenChange ScreenId

func (sc ScreenChange) Act(c *Context) {
	if !c.Bot.behaviour.ScreenExist(ScreenId(sc)) {
		panic(ScreenNotExistErr)
	}
	err := c.ChangeScreen(ScreenId(sc))
	if err != nil {
		panic(err)
	}
}

type C = Context

// Changes screen of user to the Id one.
func (c *Context) ChangeScreen(screenId ScreenId) error {
	if !c.Bot.behaviour.ScreenExist(screenId) {
		return ScreenNotExistErr
	}

	// Stop the reading by sending the nil,
	// since we change the screen and
	// current goroutine needs to be stopped.
	// if c.readingUpdate {
		// c.Updates <- nil
	// }

	// Getting the screen and changing to
	// then executing its widget.
	screen := c.Bot.behaviour.Screens[screenId]
	c.PrevScreen = c.CurScreen
	c.CurScreen = screen

	// Making the new channel for the widget.
	if c.widgetUpdates != nil {
		close(c.widgetUpdates)
	}
	c.widgetUpdates = make(chan *Update)
	if screen.Widget != nil {
		// Running the widget if the screen has one.
		go screen.Widget.Serve(c, c.widgetUpdates)
	} else {
		// Skipping updates if there is no
		// widget to handle them.
		go func() {
			for _ = range c.widgetUpdates {}
		}()
	}

	//c.Bot.Render(c.Session.Id, screen)
	//if screen.Action != nil {
		//c.run(screen.Action, c.Update)
	//}

	return nil
}
