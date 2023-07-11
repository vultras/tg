package behx

// Represents unique value to identify chats.
// In fact is simply ID of the chat.
type SessionId int64

// The type represents current state of
// user interaction per each of them.
type Session struct {
	Id SessionId
	CurrentScreenId ScreenId
	PreviousScreenId ScreenId
}

// The type represents map of sessions using
// as key.
type SessionMap map[SessionId] *Session

func (si SessionId) ToTelegram() int64 {
	return int64(si)
}

func (sm SessionMap) Add(sid SessionId) {
	sm[sid] = &Session{
		Id: sid,
	}
}

