package tg

import (
	//tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Implementing the interface provides
// ability to build your own widgets,
// aka components.
type Widget interface {
	// When the update channel is closed
	// widget MUST end its work.
	// Mostly made by looping over the
	// updates range.
	Serve(*Context, chan *Update) error
}

// Implementing the interface provides 
type DynamicWidget interface {
	MakeWidget() Widget
}

// The function that implements the Widget
// interface.
type WidgetFunc func(*Context, chan *Update) error

func (wf WidgetFunc) Serve(c *Context, updates chan *Update) error {
	return wf(c, updates)
}

