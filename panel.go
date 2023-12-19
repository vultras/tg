package tg

type Rowser interface {
	MakeRows(c *Context) []ButtonRow
}

type RowserFunc func(c *Context) []ButtonRow
func (fn RowserFunc) MakeRows(c *Context) []ButtonRow {
	return fn(c)
}

// The type represents the inline panel with
// scrollable via buttons content.
// Can be used for example to show users via SQL and offset
// or something like that.
type PanelCompo struct {
	*InlineCompo
	Rowser Rowser
}

// Transform to the panel with dynamic rows.
func (compo *MessageCompo) Panel(
	c *Context, // The context that all the buttons will get.
	rowser Rowser, // The rows generator.
) *PanelCompo {
	ret := &PanelCompo{}
	ret.InlineCompo = compo.Inline(
		NewKeyboard(
			rowser.MakeRows(c)...,
		).Inline(),
	)
	ret.Rowser = rowser
	return ret
}

func (compo *PanelCompo) Update(c *Context) {
	compo.Rows = compo.Rowser.MakeRows(c)
	compo.InlineCompo.Update(c)
}

