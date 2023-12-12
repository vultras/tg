package tg

// The type describes dynamic screen widget
// That can have multiple UI components.
type Widget interface {
	Render(*Context) UI
}

// The way to describe custom function based Widgets.
type RenderFunc func(c *Context) UI
func (fn RenderFunc) Render(c *Context) UI {
	return fn(c)
}

// The type that represents endpoint user interface
// via set of components that will work on the same screen
// in the same time.
type UI []Component

// The type describes interfaces
// needed to be implemented to be endpoint handlers.
type Component interface {
	// Optionaly component can implement the
	// Renderable interface to automaticaly be sent to the
	// user side.

	Filterer
	Server
}

