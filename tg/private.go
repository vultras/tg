package tg

import (
	"fmt"

	//tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type ContextType string

const (
	NoContextType = iota
	PrivateContextType
	GroupContextType
	ChannelContextType
)

// General context for a specific user.
// Is always the same and is not reached
// inside end function-handlers.
type context struct {
	Session *Session
	// To reach the bot abilities inside callbacks.
	Bot     *Bot
	skippedUpdates *UpdateChan
	// Current screen ID.
	screenId, prevScreenId ScreenId
}

// Goroutie function to handle each user.
func (c *Context) serve() {
	beh := c.Bot.behaviour
	if beh.Init != nil {
		c.Run(beh.Init, c.Update)
	}
	beh.Root.Serve(c)
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

func (c *Context) Sendf2(format string, v ...any) (*Message, error) {
	return c.Send(NewMessage(fmt.Sprintf(format, v...)).MD2())
}

func (c *Context) SendfHTML(format string, v ...any) (*Message, error) {
	return c.Send(NewMessage(fmt.Sprintf(format, v...)).HTML())
}

// Interface to interact with the user.
type Context struct {
	*context
	// The update that called the Context usage.
	*Update
	// Used as way to provide outer values redirection
	// into widgets and actions. It is like arguments
	// for REST API request etc.
	Arg any
	// Instead of updates as argument.
	input *UpdateChan
}

// Get the input for current widget.
// Should be used inside handlers (aka "Serve").
func (c *Context) Input() chan *Update {
	return c.input.Chan()
}

// Returns copy of current context so
// it will not affect the current one.
// But be careful because
// most of the insides uses pointers
// which are not deeply copied.
func (c *Context) Copy() *Context {
	ret := *c
	return &ret
}

func (c *Context) WithArg(v any) *Context {
	c.Arg = v
	return c
}

func (c *Context) WithUpdate(u *Update) *Context {
	c.Update = u
	return c
}

func (c *Context) WithInput(input *UpdateChan) *Context {
	c.input = input
	return c
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
	c.skippedUpdates = nil
	if screen.Widget != nil {
		c.skippedUpdates = c.RunWidget(screen.Widget, args)
	} else {
		panic("no widget defined for the screen")
	}

	return nil
}

// Run widget in background returning the new input channel for it.
func (c *Context) RunWidget(widget Widget, args ...any) *UpdateChan {
	if widget == nil {
		return nil
	}


	var arg any
	if len(args) == 1 {
		arg = args[0]
	} else if len(args) > 1 {
		arg = args
	}

	updates := NewUpdateChan()
	go func() {
		widget.Serve(
			c.Copy().
				WithInput(updates).
				WithArg(arg),
		)
		updates.Close()
	}()

	return updates
}

func (c *Context) ChangeToPrevScreen() {
	c.ChangeScreen(c.PrevScreenId())
}
