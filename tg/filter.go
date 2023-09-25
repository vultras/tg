package tg


// Implementing the interface provides way
// to know exactly what kind of updates
// the widget needs.
type Filterer interface {
	// Return true if should filter the update
	// and not send it inside the widget.
	Filter(*Update) bool
}

type FilterFunc func(*Update, MessageMap) bool
func (f FilterFunc) Filter(
	u *Update, msgs MessageMap,
) bool {
	return f(u, msgs)
}
