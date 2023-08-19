package tg

// Customized actions for the group behaviour.
type GroupAction interface {
	Act(*GroupContext)
}

// The handler function type.
type GroupActionFunc func(*GroupContext)

func (af GroupActionFunc) Act(a *GroupContext) {
	af(a)
}

type GC = GroupContext

func (c *GroupContext) SentFromSid() SessionId {
	return SessionId(c.SentFrom().ID)
}

func (a *GroupContext) SessionValue() any {
	v, _ := a.Bot.SessionValueBySid(a.SentFromSid())
	return v
}

type GroupContext struct {
	*groupContext
	*Update
}
