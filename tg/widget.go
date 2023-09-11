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
	Serve(*Context, chan *Update) error
}

// Implementing the interface provides 
type DynamicWidget interface {
	MakeWidget() Widget
}

// The function that implements the Widget
// interface.
type WidgetFunc func(*Context, chan *Update)

func (wf WidgetFunc) Serve(c *Context, updates chan *Update){
	wf(c, updates)
}

// The basic widget to provide keyboard functionality
// without implementing much.
type Page struct {
	Text string
	SubWidget Widget
	Inline *InlineKeyboard
	Reply *ReplyKeyboard
	Action Action
}

// Return new page with the specified text.
func NewPage(text string) *Page {
	ret := &Page{}
	ret.Text = text
	return ret
}

// Set the inline keyboard.
func (p *Page) WithInline(inline *InlineKeyboard) *Page {
	p.Inline = inline
	return p
}

// Set the reply keyboard.
func (p *Page) WithReply(reply *ReplyKeyboard) *Page {
	p.Reply = reply
	return p
}

// Set the action to be run before serving.
func (p *Page) WithAction(a Action) *Page {
	p.Action = a
	return p
}

// Alias to with action to simpler define actions.
func (p *Page) ActionFunc(fn ActionFunc) *Page {
	return p.WithAction(fn)
}

// Set the sub widget that will get the skipped
// updates.
func (p *Page) WithSub(sub Widget) *Page {
	p.SubWidget = sub
	return p
}

func (p *Page) Serve(
	c *Context, updates chan *Update,
) error {
	msgs, err := c.Render(p)
	if err != nil {
		return err
	}

	// The inline message is always returned
	// and the reply one is useless in our case.
	inlineMsg := msgs[0]

	if p.Action != nil {
		c.run(p.Action, c.Update)
	}
	var subUpdates chan *Update
	if p.SubWidget != nil {
		subUpdates = make(chan *Update)
		go p.SubWidget.Serve(c, subUpdates)
		defer close(subUpdates)
	}
	for u := range updates {
		var act Action
		if u.Message != nil {
			text := u.Message.Text
			kbd := p.Reply
			if kbd == nil {
				if subUpdates != nil {
					subUpdates <- u
				}
				continue
			}
			btns := kbd.ButtonMap()
			btn, ok := btns[text]
			if !ok {
				if u.Message.Location != nil {
					for _, b := range btns {
						if b.SendLocation {
							btn = b
							ok = true
						}
					}
				} else if subUpdates != nil {
					subUpdates <- u
				}
			}
			if btn != nil {
				act = btn.Action
			} else if kbd.Action != nil {
				act = kbd.Action
			}
		} else if u.CallbackQuery != nil {
			if u.CallbackQuery.Message.MessageID != inlineMsg.MessageID {
				if subUpdates != nil {
					subUpdates <- u
				}
				continue
			}
			cb := tgbotapi.NewCallback(
				u.CallbackQuery.ID,
				u.CallbackQuery.Data,
			)
			data := u.CallbackQuery.Data

			_, err := c.Bot.Api.Request(cb)
			if err != nil {
				return err
			}
			kbd := p.Inline
			if kbd == nil {
				if subUpdates != nil {
					subUpdates <- u
				}
				continue
			}

			btns := kbd.ButtonMap()
			btn, ok := btns[data]
			if !ok {
				if subUpdates != nil {
					subUpdates <- u
				}
				continue
			}
			if btn != nil {
				act = btn.Action
			} else if kbd.Action != nil {
				act = kbd.Action
			}
		}
		if act != nil {
			c.run(act, u)
		} 
	}
	return nil
}

func (s *Page) Render(
	sid SessionId, bot *Bot,
) ([]*SendConfig, error) {
	cid := sid.ToApi()
	reply := s.Reply
	inline := s.Inline
	ret := []*SendConfig{}
	var txt string
	// Screen text and inline keyboard.
	if s.Text != "" {
		txt = s.Text
	} else if inline != nil {
		// Default to send the keyboard.
		txt = ">"
	}
	if txt != "" {
		msgConfig := tgbotapi.NewMessage(cid, txt)
		if inline != nil {
			msgConfig.ReplyMarkup = inline.ToApi()
		} else if reply != nil {
			msgConfig.ReplyMarkup = reply.ToApi()
			ret = append(ret, &SendConfig{Message: &msgConfig})
			return ret, nil
		} else {
			msgConfig.ReplyMarkup = NewReply().
				WithRemove(true).
				ToApi()
			ret = append(ret, &SendConfig{Message: &msgConfig})
			return ret, nil
		}
		ret = append(ret, &SendConfig{Message: &msgConfig})
	}

	// Screen text and reply keyboard.
	if reply != nil {
		msgConfig := tgbotapi.NewMessage(cid, ">")
		msgConfig.ReplyMarkup = reply.ToApi()
		ret = append(ret, &SendConfig{
			Message: &msgConfig,
		})
	} else {
		// Removing keyboard if there is none.
		msgConfig := tgbotapi.NewMessage(cid, ">")
		msgConfig.ReplyMarkup = NewReply().
			WithRemove(true).
			ToApi()
			ret = append(ret, &SendConfig{Message: &msgConfig})
	}

	return ret, nil
}

