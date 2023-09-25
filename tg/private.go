package tg

import (
	"fmt"

	//tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	//"path"
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
	path, prevPath Path
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

func (c *Context) Path() Path {
	return c.path
}

func (c *Context) PrevPath() Path {
	return c.prevPath
}

func (c *Context) Run(a Action, u *Update) {
	if a != nil {
		a.Act(c.Copy().WithUpdate(u))
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


type C = Context

// Changes screen of user to the Id one.
func (c *Context) Go(pth Path, args ...any) error {
	// Getting the screen and changing to
	// then executing its widget.
	if !pth.IsAbs() {
		pth = (c.Path() + "/" + pth).Clean()
	}

	if !c.PathExist(pth) {
		return ScreenNotExistErr
	}
	c.prevPath = c.path
	c.path = pth

	// Stopping the current widget.
	screen := c.Bot.behaviour.Screens[pth]
	c.skippedUpdates.Close()
	if screen.Widget != nil {
		c.skippedUpdates = c.RunWidget(screen.Widget, args...)
	} else {
		panic("no widget defined for the screen")
	}

	return nil
}

func (c *Context) PathExist(pth Path) bool {
	return c.Bot.behaviour.PathExist(pth)
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
		// To let widgets finish themselves before
		// the channel is closed.
		updates.Close()
	}()

	return updates
}

// Simple way to read strings for widgets.
func (c *Context) ReadString(pref string, args ...any) string {
	var text string
	c.Sendf(pref, args...)
	for u := range c.Input() {
		if u.Message == nil {
			continue
		}
		text = u.Message.Text
		break
	}
	return text
}

// Change screen to the previous.
// To get to the parent screen use GoUp.
func (c *Context) GoPrev() {
	pth := c.PrevPath()
	if pth == "" {
		c.Go("/")
	}
	c.Go(pth)
}
