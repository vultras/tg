package tg

// Implementing the interface provides
// the way to define how to handle updates.
type Server interface {
	Serve(*Context)
}

