package tx

// Represents unique value to identify chats.
// In fact is simply ID of the chat.
type SessionId int64

// The type represents current state of
// user interaction per each of them.
type Session struct {
	Id SessionId
	// Current screen identifier.
	CurrentScreenId ScreenId
	// ID of the previous screen.
	PreviousScreenId ScreenId
	// The currently showed on display keyboard inside Action.
	KeyboardId KeyboardId
	V          any
}

// The type represents map of sessions using
// as key.
type SessionMap map[SessionId]*Session

// Session information for a group.
type GroupSession struct {
	Id SessionId
	// Information for each user in the group.
	V map[SessionId]any
}

// Map for every user in every chat sessions.
type GroupSessionMap map[SessionId]*GroupSession

// Return new empty session with specified user ID.
func NewSession(id SessionId) *Session {
	return &Session{
		Id: id,
		V:  make(map[string]any),
	}
}

// Returns new empty group session with specified group and user IDs.
func NewGroupSession(id SessionId) *GroupSession {
	return &GroupSession{
		Id: id,
		V:  make(map[SessionId]any),
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
