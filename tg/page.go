package tg

import (
	//tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// The basic widget to provide keyboard functionality
// without implementing much.
type Page struct {
	Text string
	SubWidget Widget
	Inline *InlineKeyboardWidget
	Action Action
	Reply *ReplyKeyboardWidget
}

// Return new page with the specified text.
func NewPage(text string) *Page {
	ret := &Page{}
	ret.Text = text
	return ret
}

// Set the inline keyboard.
func (p *Page) WithInline(inline *InlineKeyboardWidget) *Page {
	p.Inline = inline
	return p
}

// Set the reply keyboard.
func (p *Page) WithReply(reply *ReplyKeyboardWidget) *Page {
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


func (p *Page) Render(
	sid SessionId, bot *Bot,
) ([]*SendConfig) {
	reply := p.Reply
	inline := p.Inline

	ret := []*SendConfig{}

	if p.Text != "" {
		cfg := NewMessage(p.Text).SendConfig(sid, bot).
			WithName("page/text")
		ret = append(ret, cfg)
	}
	if inline != nil {
		cfg := inline.SendConfig(sid, bot).
			WithName("page/inline")
		ret = append(ret, cfg)
	}
	if p.Reply != nil {
		cfg := reply.SendConfig(sid, bot).
			WithName("page/reply")
		ret = append(ret, cfg)
	}

	return ret
}

func (p *Page) Filter(
	u *Update, msgs MessageMap,
) bool {
	return false
}

func (p *Page) Serve(
	c *Context, updates *UpdateChan,
) {
	msgs, _ := c.Render(p)
	inlineMsg := msgs["page/inline"]
	if p.Action != nil {
		c.Run(p.Action, c.Update)
	}

	subUpdates := c.RunWidgetBg(p.SubWidget)
	defer subUpdates.Close()

	inlineUpdates := c.RunWidgetBg(p.Inline)
	defer inlineUpdates.Close()

	replyUpdates := c.RunWidgetBg(p.Reply)
	defer replyUpdates.Close()

	subFilter, subFilterOk := p.SubWidget.(Filterer)
	for u := range updates.Chan() {
		switch {
		case !p.Inline.Filter(u, MessageMap{"": inlineMsg}) :
			inlineUpdates.Send(u)
		case !p.Reply.Filter(u, msgs) :
			replyUpdates.Send(u )
		case p.SubWidget != nil :
			if subFilterOk {
				if !subFilter.Filter(u, msgs) {
					subUpdates.Send(u)
				}
			} else {
				subUpdates.Send(u)
			}
		default:
		}
	}
}


