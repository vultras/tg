package tg

// The type to descsribe one line reading widget.
type UpdateRead struct {
	Pre Action
	Filterer Filterer
	Post Widget
}

func (rd *UpdateRead) Filter(u *Update, msgs MessageMap) bool {
	if rd.Filterer != nil {
		return rd.Filterer.Filter(u, msgs)
	}

	return false
}

// Returns new empty update reader.
func NewUpdateRead(filter Filterer, post Widget) *UpdateRead {
	ret := &UpdateRead{}
	ret.Filterer = filter
	ret.Post = post
	return ret
}

func (rd *UpdateRead) WithPre(a Action) *UpdateRead {
	rd.Pre = a
	return rd
}

func NewTextMessageRead(pre Action, post Widget) *UpdateRead {
	ret := NewUpdateRead(
		FilterFunc(func(u *Update, _ MessageMap) bool {
			return u.Message == nil
		}),
		post,
	).WithPre(pre)
	return ret
}

func (rd *UpdateRead) Serve(c *Context) {
	c.Run(rd.Pre, c.Update)
	for u := range c.Input() {
		if rd.Filter(u, nil) {
			continue
		}
		c.RunWidget(rd.Post, u)
		break
	}
}


