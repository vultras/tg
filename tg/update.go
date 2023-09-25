package tg

// The type represents general update channel.
type UpdateChan struct {
	chn chan *Update
}

// Return new update channel.
func NewUpdateChan() *UpdateChan {
	ret := &UpdateChan{}
	ret.chn = make(chan *Update)
	return ret
}


func (updates *UpdateChan) Chan() chan *Update {
	return updates.chn
}

// Send an update to the channel.
// Returns true if the update was sent.
func (updates *UpdateChan) Send(u *Update) bool {
	if updates == nil || updates.chn == nil {
		return false
	}
	updates.chn <- u
	return true
}

// Read an update from the channel.
func (updates *UpdateChan) Read() *Update {
	if updates == nil || updates.chn == nil {
		return nil
	}
	return <-updates.chn
}

// Returns true if the channel is closed.
func (updates *UpdateChan) Closed() bool {
	return updates==nil || updates.chn == nil
}

// Close the channel. Used in defers.
func (updates *UpdateChan) Close() {
	if updates == nil || updates.chn == nil {
		return
	}
	close(updates.chn)
	updates.chn = nil
}

