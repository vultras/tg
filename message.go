package tg

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	//"strings"
	re "regexp"
	"fmt"
)
type Message = tgbotapi.Message

// Simple text message component type.
type MessageCompo struct {
	Message *Message
	ParseMode string
	Text string
}

var (
	escapeRe = re.MustCompile(`([_*\[\]()~`+"`"+`>#+-=|{}.!])`)
	NewRawMessage = tgbotapi.NewMessage
)

// Escape special characters in Markdown 2 and return the
// resulting string.
func Escape2(str string) string {
	return string(escapeRe.ReplaceAll([]byte(str), []byte("\\$1")))
}

func (compo *MessageCompo) Update(c *Context) {
	edit := tgbotapi.NewEditMessageText(
		c.Session.Id.ToApi(),
		compo.Message.MessageID,
		compo.Text,
	)
	msg, _ := c.Bot.Api.Send(edit)
	compo.Message = &msg
}

// Is only implemented to make it sendable and so we can put it
// return of rendering functions.
func (compo *MessageCompo) SetMessage(msg *Message) {
	compo.Message = msg
}

// Return new message with the specified text.
func NewMessage(format string, v ...any) *MessageCompo {
	ret := &MessageCompo{}
	ret.Text = fmt.Sprintf(format, v...)
	ret.ParseMode = tgbotapi.ModeMarkdown
	return ret
}

// Return message with the specified parse mode.
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

// Transform the message component into the location one.
func (msg *MessageCompo) Location(
	lat, long float64,
) *LocationCompo {
	ret := &LocationCompo{
		MessageCompo: msg,
		Location: Location{
			Latitude: lat,
			Longitude: long,
		},
	}
	return ret
}

// Implementing the Sendable interface.
func (config *MessageCompo) SendConfig(
	sid SessionId, bot *Bot,
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

	msg := tgbotapi.NewMessage(sid.ToApi(), text)
	ret.Message = &msg
	ret.Message.ParseMode = config.ParseMode

	return &ret
}

// Empty serving to use messages in rendering.
func (compo *MessageCompo) Serve(c *Context) {}

// Filter that skips everything. Messages cannot do anything with updates.
func (compo *MessageCompo) Filter(_ *Update) bool {
	// Skip everything
	return true
}
