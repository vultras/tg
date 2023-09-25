package tg

type UIs []UI

// The type describes dynamic screen widget.
type Widget interface {
	UIs(*Context) UIs
}

// The way to describe custom function based Widgets.
type WidgetFunc func(c *Context) UIs
func (fn WidgetFunc) UIs(c *Context) UIs {
	return fn(c)
}

// The type describes interfaces
// needed to be implemented to be endpoint handlers.
type UI interface {
	Renderable

	SetMessage(*Message)
	GetMessage() *Message
	Filterer

	Server
}

type UiFunc func()

// The type to embed into potential components.
// Implements empty versions of interfaces
// and contains 
type Compo struct{
	*Message
}

// Defalut setting message 
func (compo Compo) SetMessage(msg *Message) { compo.Message = msg }
func (compo Compo) GetMessage() *Message { return compo.Message }
// Default non filtering filter. Always returns false.
func (compo Compo) Filter(_ *Update, _ *Message) bool {return false}

