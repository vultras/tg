package tx

import (
	apix "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Update = apix.Update

// Represents unique value to identify chats.
// In fact is simply ID of the chat.
type SessionId int64

// The type represents current state of
// user interaction per each of them.
type Session struct {
	// Unique identifier for the session, Telegram chat's ID.
	Id SessionId
	// Current screen identifier.
	CurrentScreenId ScreenId
	// ID of the previous screen.
	PreviousScreenId ScreenId
	// The currently showed on display keyboard inside Action.
	KeyboardId KeyboardId

	// Is true if currently reading the Update.
	readingUpdate bool

	// Custom data for each user.
	V map[string]any
}

// The type represents map of sessions using
// as key.
type SessionMap map[SessionId]*Session

// Return new empty session with
func NewSession(id SessionId) *Session {
	return &Session{
		Id: id,
		V:  make(map[string]any),
	}
}

// Changes screen of user to the Id one for the session.
func (c *Session) ChangeScreen(screenId ScreenId) {
	c.PreviousScreenId = c.CurrentScreenId
	c.CurrentScreenId = screenId
}

// Convert the SessionId to Telegram API's type.
func (si SessionId) ToTelegram() int64 {
	return int64(si)
}

// Add new empty session by it's ID.
func (sm SessionMap) Add(sid SessionId) {
	sm[sid] = NewSession(sid)
}
