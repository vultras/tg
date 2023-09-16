package tg

import (
	"fmt"

	//tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type context struct {
	Session *Session
	// To reach the bot abilities inside callbacks.
	Bot     *Bot
	skippedUpdates *UpdateChan
	// Current screen ID.
	screenId, prevScreenId ScreenId
}


// The type represents way to interact with user in
// handling functions. Is provided to Act() function always.

// Goroutie function to handle each user.
func (c *Context) Serve(updates *UpdateChan) {
	beh := c.Bot.behaviour
	if beh.Init != nil {
		c.Run(beh.Init, c.Update)
	}
	beh.Root.Serve(c, updates)
}


func (c *context) run(a Action, u *Update) {
	a.Act(&Context{context: c, Update:  u})
}

func (c *Context) ScreenId() ScreenId {
	return c.screenId
}

func (c *Context) PrevScreenId() ScreenId {
	return c.prevScreenId
}

func (c *Context) Run(a Action, u *Update) {
	if a != nil {
		a.Act(&Context{context: c.context, Update: u})
	}
}

// Only for the root widget usage.
// Skip the update sending it down to
// the underlying widget.
func (c *Context) Skip(u *Update) {
	c.skippedUpdates.Send(u)
}

// Renders the Renedrable object to the side of client
// and returns the messages it sent.
func (c *Context) Render(v Renderable) (MessageMap, error) {
	return c.Bot.Render(c.Session.Id, v)
}

// Sends to the Sendable object.
func (c *Context) Send(v Sendable) (*Message, error) {
	return c.Bot.Send(c.Session.Id, v)
}

// Sends the formatted with fmt.Sprintf message to the user.
func (c *Context) Sendf(format string, v ...any) (*Message, error) {
	return c.Send(NewMessage(fmt.Sprintf(format, v...)))
}

// Interface to interact with the user.
type Context struct {
	*context
	// The update that called the Context usage.
	*Update
	// Used as way to provide outer values redirection
	// into widgets and actions. It is like arguments
	// for REST API request etc.
	Args []any
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
func (c *Context) ChangeScreen(screenId ScreenId, args ...any) error {
	if !c.Bot.behaviour.ScreenExist(screenId) {
		return ScreenNotExistErr
	}

	// Getting the screen and changing to
	// then executing its widget.
	screen := c.Bot.behaviour.Screens[screenId]
	c.prevScreenId = c.screenId
	c.screenId = screenId

	// Stopping the current widget.
	c.skippedUpdates.Close()
	// Making channel for the new widget.
	c.skippedUpdates = NewUpdateChan()
	if screen.Widget != nil {
		// Running the widget if the screen has one.
		go func() {
			updates := c.skippedUpdates
			screen.Widget.Serve(&Context{
				context: c.context,
				Update: c.Update,
				Args: args,
			}, updates)
			updates.Close()
		}()
	} else {
		panic("no widget defined for the screen")
	}

	return nil
}

func (c *Context) ChangeToPrevScreen() {
	c.ChangeScreen(c.PrevScreenId())
}
