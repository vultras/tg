package tg

// The way to determine where the context is
// related to.
type SessionScope uint8
const (
	NoSessionScope SessionScope = iota
	PrivateSessionScope
	GroupSessionScope
	ChannelSessionScope
)

// Represents unique value to identify chats.
// In fact is simply ID of the chat.
type SessionId int64

// Convert the SessionId to Telegram API's type.
func (si SessionId) ToApi() int64 {
	return int64(si)
}

// The type represents current state of
// user interaction per each of them.
type Session struct {
	// Id of the chat of the user.
	Id SessionId
	Scope SessionScope
	// Custom value for each user.
	Data  any
}

// Return new empty session with specified user ID.
func NewSession(id SessionId, scope SessionScope) *Session {
	return &Session{
		Id: id,
		Scope: scope,
	}
}

// The type represents map of sessions using
// as key.
type SessionMap map[SessionId]*Session

// Add new empty session by it's ID.
func (sm SessionMap) Add(sid SessionId, scope SessionScope) *Session {
	ret := NewSession(sid, scope)
	sm[sid] = ret
	return ret
}

// Session information for a group.
type GroupSession struct {
	Id SessionId
	// Information for each user in the group.
	Data any
}

// Returns new empty group session with specified group and user IDs.
func NewGroupSession(id SessionId) *GroupSession {
	return &GroupSession{
		Id: id,
	}
}

// Map for every group the bot is in.
type GroupSessionMap map[SessionId]*GroupSession

func (sm GroupSessionMap) Add(sid SessionId) {
	sm[sid] = NewGroupSession(sid)
}
