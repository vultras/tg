package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	//"strings"
)
type Message = tgbotapi.Message

// Simple text message type.
type MessageCompo struct {
	Message *Message
	ParseMode string
	Text string
}

func (compo *MessageCompo) SetMessage(msg *Message) {
	compo.Message = msg
}

// Return new message with the specified text.
func NewMessage(text string) *MessageCompo {
	ret := &MessageCompo{}
	ret.Text = text
	ret.ParseMode = tgbotapi.ModeMarkdown
	return ret
}

func (msg *MessageCompo) withParseMode(mode string) *MessageCompo {
	msg.ParseMode = mode
	return msg
}

// Set the default Markdown parsing mode.
func (msg *MessageCompo) MD() *MessageCompo {
	return msg.withParseMode(tgbotapi.ModeMarkdown)
}

// Set the Markdown 2 parsing mode.
func (msg *MessageCompo) MD2() *MessageCompo {
	return msg.withParseMode(tgbotapi.ModeMarkdownV2)
}


// Set the HTML parsing mode.
func (msg *MessageCompo) HTML() *MessageCompo {
	return msg.withParseMode(tgbotapi.ModeHTML)
}

// Transform the message component into one with reply keyboard.
func (msg *MessageCompo) Inline(inline *Inline) *InlineCompo {
	return &InlineCompo{
		Inline: inline,
		MessageCompo: msg,
	}
}

// Transform the message component into one with reply keyboard.
func (msg *MessageCompo) Reply(reply *Reply) *ReplyCompo {
	return &ReplyCompo{
		Reply: reply,
		MessageCompo: msg,
	}
}

func (config *MessageCompo) SendConfig(
	c *Context,
) (*SendConfig) {
	var (
		ret SendConfig
		text string
	)

	if config.Text == "" {
		text = ">"
	} else {
		text = config.Text
	}

	//text = strings.ReplaceAll(text, "-", "\\-")

	msg := tgbotapi.NewMessage(c.Session.Id.ToApi(), text)
	ret.Message = &msg
	ret.Message.ParseMode = config.ParseMode

	return &ret
}

// Empty serving to use messages in rendering.
func (compo *MessageCompo) Serve(c *Context) {
}

func (compo *MessageCompo) Filter(_ *Update) bool {
	// Skip everything
	return true
}
