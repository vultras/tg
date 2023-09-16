package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// Implementing the interface provides
// ability to build your own widgets,
// aka components.
type Widget interface {
	// When the update channel is closed
	// widget MUST end its work.
	// Mostly made by looping over the
	// updates range.
	Serve(*Context, *UpdateChan)
}

// Needs implementation.
// Behaviour can be the root widget or something like
// that.
type RootWidget interface {
	Widget
}

// Implementing the interface provides way
// to know exactly what kind of updates
// the widget needs.
type Filterer interface {
	// Return true if should filter the update
	// and not send it inside the widget.
	Filter(*Update, MessageMap) bool
}

// The type represents general update channel.
type UpdateChan struct {
	chn chan *Update
}

// Return new update channel.
func NewUpdateChan() *UpdateChan {
	ret := &UpdateChan{}
	ret.chn = make(chan *Update)
	return ret
}

func (updates *UpdateChan) Chan() chan *Update {
	return updates.chn
}

// Send an update to the channel.
func (updates *UpdateChan) Send(u *Update) {
	if updates != nil && updates.chn == nil {
		return
	}
	updates.chn <- u
}

// Read an update from the channel.
func (updates *UpdateChan) Read() *Update {
	if updates == nil || updates.chn == nil {
		return nil
	}
	return <-updates.chn
}

// Returns true if the channel is closed.
func (updates *UpdateChan) Closed() bool {
	return updates.chn == nil
}

// Close the channel. Used in defers.
func (updates *UpdateChan) Close() {
	if updates == nil || updates.chn == nil {
		return
	}
	close(updates.chn)
	updates.chn = nil
}

func (c *Context) RunWidgetBg(widget Widget) *UpdateChan {
	if widget == nil {
		return nil
	}

	updates := NewUpdateChan()
	go widget.Serve(c, updates)

	return updates
}

// Implementing the interface provides 
type DynamicWidget interface {
	MakeWidget() Widget
}

// The function that implements the Widget
// interface.
type WidgetFunc func(*Context, *UpdateChan)

func (wf WidgetFunc) Serve(c *Context, updates *UpdateChan) {
	wf(c, updates)
}

func (wf WidgetFunc) Filter(
	u *Update,
	msgs ...*Message,
) bool {
	return false
}

// The type implements message with an inline keyboard.
type InlineKeyboardWidget struct {
	Text string
	*InlineKeyboard
}

// The type implements dynamic inline keyboard widget.
// Aka message with inline keyboard.
func NewInlineKeyboardWidget(
	text string,
	inline *InlineKeyboard,
) *InlineKeyboardWidget {
	ret := &InlineKeyboardWidget{}
	ret.Text = text
	ret.InlineKeyboard = inline
	return ret
}
func (widget *InlineKeyboardWidget) SendConfig(
	sid SessionId,
	bot *Bot,
) (*SendConfig) {
	var text string
	if widget.Text != "" {
		text = widget.Text
	} else {
		text = ">"
	}

	msgConfig := tgbotapi.NewMessage(sid.ToApi(), text)
	msgConfig.ReplyMarkup = widget.ToApi()

	ret := &SendConfig{}
	ret.Message = &msgConfig
	return ret
}

func (widget *InlineKeyboardWidget) Serve(
	c *Context,
	updates *UpdateChan,
) {
	for u := range updates.Chan() {
		var act Action
		if u.CallbackQuery == nil {
			continue
		}
		cb := tgbotapi.NewCallback(
			u.CallbackQuery.ID,
			u.CallbackQuery.Data,
		)
		data := u.CallbackQuery.Data

		_, err := c.Bot.Api.Request(cb)
		if err != nil {
			//return err
			continue
		}

		btns := widget.ButtonMap()
		btn, ok := btns[data]
		if !ok {
			continue
		}
		if btn != nil {
			act = btn.Action
		} else if widget.Action != nil {
			act = widget.Action
		}
		c.Run(act, u)
	}
}

func (widget *InlineKeyboardWidget) Filter(
	u *Update,
	msgs MessageMap,
) bool {
	if widget == nil {
		return true
	}
	if u.CallbackQuery == nil || len(msgs) < 1 {
		return true
	}

	inlineMsg, inlineOk := msgs[""]
	if inlineOk {
		if u.CallbackQuery.Message.MessageID != 
				inlineMsg.MessageID {
			return true
		}
	}

	return false
}

// The type implements dynamic reply keyboard widget.
type ReplyKeyboardWidget struct {
	Text string
	*ReplyKeyboard
}

// Returns new empty reply keyboard widget.
func NewReplyKeyboardWidget(
	text string,
	reply *ReplyKeyboard,
) *ReplyKeyboardWidget {
	ret := &ReplyKeyboardWidget{}
	ret.Text = text
	ret.ReplyKeyboard = reply
	return ret
}

func (widget *ReplyKeyboardWidget) SendConfig(
	sid SessionId,
	bot *Bot,
) (*SendConfig) {
	var text string
	if widget.Text != "" {
		text = widget.Text
	} else {
		text = ">"
	}

	msgConfig := tgbotapi.NewMessage(sid.ToApi(), text)
	msgConfig.ReplyMarkup = widget.ToApi()

	ret := &SendConfig{}
	ret.Message = &msgConfig
	return ret
}

func (widget *ReplyKeyboardWidget) Filter(
	u *Update,
	msgs MessageMap,
) bool {
	if widget == nil {
		return true
	}
	if u.Message == nil {
		return true
	}
	_, ok := widget.ButtonMap()[u.Message.Text]
	if !ok {
		return true
	}
	return false
}

func (widget *ReplyKeyboardWidget) Serve(
	c *Context,
	updates *UpdateChan,
) {
	for u := range updates.Chan() {
		var btn *Button
		text := u.Message.Text
		btns := widget.ButtonMap()
		btn, ok := btns[text]
		if !ok {
			if u.Message.Location != nil {
				for _, b := range btns {
					if b.SendLocation {
						btn = b
						ok = true
					}
				}
			}
		}

		if btn != nil {
			c.Run(btn.Action, u)
		}
	}
}

