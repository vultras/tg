package tg

// Implementing the interface lets the
// value to be sent.
type Sendable interface {
	Send(SessionId, *Bot) (*Message, error)
}

type Renderable interface {
	Render(SessionId, *Bot) ([]*Message, error)
}
