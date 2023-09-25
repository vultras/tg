package tg

import (
	//tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	//"fmt"
)

type Maker[V any] interface {
	Make(*Context) V
}

type MakeFunc[V any] func(*Context) V
func (fn MakeFunc[V]) Make(c *Context) V {
	return fn(c)
}

type ArgMap = map[string] any
type ArgSlice = []any
type ArgList[V any] []V

// Implementing the interface provides
// ability to build your own widgets,
// aka components.
type Widget interface {
	// When the update channel is closed
	// widget MUST end its work.
	// Mostly made by looping over the
	// updates range.
	Serve(*Context)
}

type DynamicWidget[W Widget] interface {
	Maker[W]
}

// Implementing the interface provides ability to
// be used as the root widget for contexts.
type RootWidget interface {
	Widget
	SetSub(Widget)
}

// Implementing the interface provides way
// to know exactly what kind of updates
// the widget needs.
type Filterer interface {
	// Return true if should filter the update
	// and not send it inside the widget.
	Filter(*Update, MessageMap) bool
}

type FilterFunc func(*Update, MessageMap) bool
func (f FilterFunc) Filter(
	u *Update, msgs MessageMap,
) bool {
	return f(u, msgs)
}

// General type function for faster typing.
type Func func(*Context)
func (f Func) Act(c *Context) {
	f(c)
}
func (f Func) Serve(c *Context) {
	f(c)
}


// The function that implements the Widget
// interface.
type WidgetFunc func(*Context)

func (wf WidgetFunc) Serve(c *Context) {
	wf(c)
}


