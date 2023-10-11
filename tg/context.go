package tg

import (
	"fmt"

	//tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	//"path"
)

// General context for a specific user.
// Is always the same and is not reached
// inside end function-handlers.
type context struct {
	Session *Session
	// To reach the bot abilities inside callbacks.
	Bot     *Bot
	// Costum status for currently running context.
	Status any
	Type ContextType
	updates *UpdateChan
	skippedUpdates *UpdateChan
	// Current screen ID.
	pathHistory []Path
	//path, prevPath Path
}

// Interface to interact with the user.
type Context struct {
	*context
	// The update that called the Context usage.
	*Update
	// Used as way to provide outer values redirection
	// into widgets and actions. It is like arguments
	// for REST API request etc.
	arg any
	// Instead of updates as argument.
	input *UpdateChan
}

// General type function to define actions, single component widgets
// and components themselves.
type Func func(*Context)
func (f Func) Act(c *Context) {
	f(c)
}
func (f Func) Serve(c *Context) {
	f(c)
}
func(f Func) Filter(_ *Update) bool {
	return false
}
func (f Func) Render(_ *Context) UI {
	return UI{
		f,
	}
}

type ContextType uint8
const (
	NoContextType ContextType = iota
	WidgetContextType
	ActionContextType
)

// Goroutie function to handle each user.
func (c *Context) serve() {
	beh := c.Bot.behaviour
	c.Run(beh.Init)
	beh.Root.Serve(c)
}

func (c *Context) Path() Path {
	ln := len(c.pathHistory)
	if ln == 0 {
		return ""
	}
	return c.pathHistory[ln-1]
}

func (c *Context) Arg() any {
	return c.arg
}

func (c *Context) Run(a Action) {
	if a != nil {
		a.Act(c)
	}
}

// Only for the root widget usage.
// Skip the update sending it down to
// the underlying widget.
func (c *Context) Skip(u *Update) {
	c.skippedUpdates.Send(u)
}

// Sends to the Sendable object.
func (c *Context) Send(v Sendable) (*Message, error) {
	config := v.SendConfig(c.Session.Id, c.Bot)
	if config.Error != nil {
		return nil, config.Error
	}

	msg, err := c.Bot.Api.Send(config.ToApi())
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

// Sends the formatted with fmt.Sprintf message to the user
// using default Markdown parsing format.
func (c *Context) Sendf(format string, v ...any) (*Message, error) {
	return c.Send(NewMessage(fmt.Sprintf(format, v...)))
}

// Same as Sendf but uses Markdown 2 format for parsing.
func (c *Context) Sendf2(format string, v ...any) (*Message, error) {
	return c.Send(NewMessage(fmt.Sprintf(format, v...)).MD2())
}

// Same as Sendf but uses HTML format for parsing.
func (c *Context) SendfHTML(format string, v ...any) (*Message, error) {
	return c.Send(NewMessage(fmt.Sprintf(format, v...)).HTML())
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
	c = c.Copy()
	c.arg = v
	return c
}

func (c *Context) WithUpdate(u *Update) *Context {
	c = c.Copy()
	c.Update = u
	return c
}

func (c *Context) WithInput(input *UpdateChan) *Context {
	c = c.Copy()
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

func (c *Context) History() []Path {
	return c.pathHistory
}

// Changes screen of user to the Id one.
func (c *Context) Go(pth Path, args ...any) {
	var back bool
	if pth == "-" {
		ln := len(c.pathHistory)
		if ln <= 1 {
			pth = "/"
		} else {
			pth = c.pathHistory[ln-2]
			c.pathHistory = c.pathHistory[:ln-1]
			back = true
		}
	}
	// Getting the screen and changing to
	// then executing its widget.
	if !pth.IsAbs() {
		pth = (c.Path() + "/" + pth).Clean()
	}

	if !c.PathExist(pth) {
		panic(ScreenNotExistErr)
	}

	if !back && c.Path() != pth {
		c.pathHistory = append(c.pathHistory, pth)
	}

	// Stopping the current widget.
	screen := c.Bot.behaviour.Screens[pth]
	c.skippedUpdates.Close()
	if screen.Widget != nil {
		c.skippedUpdates = c.RunWidget(screen.Widget, args...)
	} else {
		panic("no widget defined for the screen")
	}
}

func (c *Context) PathExist(pth Path) bool {
	return c.Bot.behaviour.PathExist(pth)
}

func (c *Context) makeArg(args []any) any {
	var arg any
	if len(args) == 1 {
		arg = args[0]
	} else if len(args) > 1 {
		arg = args
	}
	return arg
}

func (c *Context) RunCompo(compo Component, args ...any) *UpdateChan {
	if compo == nil {
		return nil
	}
	s, ok := compo.(Sendable)
	if ok {
		msg, err := c.Send(s)
		if err != nil {
			panic("could not send the message")
		}
		s.SetMessage(msg)
	}
	updates := NewUpdateChan()
	go func() {
		compo.Serve(
			c.WithInput(updates).
				WithArg(c.makeArg(args)),
		)
		// To let widgets finish themselves before
		// the channel is closed and close it by themselves.
		updates.Close()
	}()
	return updates
}

// Run widget in background returning the new input channel for it.
func (c *Context) RunWidget(widget Widget, args ...any) *UpdateChan {
	if widget == nil {
		return nil
	}

	pth := c.Path()
	compos := widget.Render(c.WithArg(c.makeArg(args)))
	// Leave if changed path.
	if compos == nil || pth != c.Path() {
		return nil
	}
	chns := make([]*UpdateChan, len(compos))
	for i, compo := range compos {
		chns[i] = c.RunCompo(compo, args...)
	}

	ret := NewUpdateChan()
	go func() {
		ln := len(compos)
		UPDATE:
		for u := range ret.Chan() {
			if u == nil {
				break
			}
			cnt := 0
			for i, compo := range compos {
				chn := chns[i]
				if chn.Closed() {
					cnt++
					continue
				}
				if !compo.Filter(u) {
					chn.Send(u)
					continue UPDATE
				}
			}
			if cnt == ln {
				break
			}
		}
		ret.Close()
		for _, chn := range chns {
			chn.Close()
		}
	}()

	return ret
}

// Simple way to read strings for widgets.
func (c *Context) ReadString(pref string, args ...any) string {
	var text string
	if pref != "" {
		c.Sendf(pref, args...)
	}
	for u := range c.Input() {
		if u == nil {
			break
		}
		if u.Message == nil {
			continue
		}
		text = u.Message.Text
		break
	}
	return text
}

