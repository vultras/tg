package tx

// Represents unique value to identify chats.
// In fact is simply ID of the chat.
type SessionId int64

// Convert the SessionId to Telegram API's type.
func (si SessionId) ToTelegram() int64 {
	return int64(si)
}

// The type represents current state of
// user interaction per each of them.
type Session struct {
	// Id of the chat of the user.
	Id SessionId
	// Custom value for each user.
	V any
}

// Return new empty session with specified user ID.
func NewSession(id SessionId) *Session {
	return &Session{
		Id: id,
		V:  make(map[string]any),
	}
}

// The type represents map of sessions using
// as key.
type SessionMap map[SessionId]*Session

// Add new empty session by it's ID.
func (sm SessionMap) Add(sid SessionId) {
	sm[sid] = NewSession(sid)
}

// Session information for a group.
type GroupSession struct {
	Id SessionId
	// Information for each user in the group.
	V map[SessionId]any
}

// Returns new empty group session with specified group and user IDs.
func NewGroupSession(id SessionId) *GroupSession {
	return &GroupSession{
		Id: id,
		V:  make(map[SessionId]any),
	}
}

// Map for every group the bot is in.
type GroupSessionMap map[SessionId]*GroupSession

func (sm GroupSessionMap) Add(sid SessionId) {
	sm[sid] = NewGroupSession(sid)
}
