package tg

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type FileId string

type Update struct {
	*tgbotapi.Update
	c *Context
}

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
	defer recover()
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
	chn := updates.chn
	updates.chn = nil
	close(chn)
}

func (u *Update) HasDocument() bool {
	return u.Message != nil && u.Message.Document != nil
}

func (u *Update) DocumentId() FileId {
	return FileId(u.Update.Message.Document.FileID)
}

func (u *Update) HasPhotos() bool {
	return u.Message != nil && u.Message.Photo != nil &&
		len(u.Message.Photo) != 0
}

func (u *Update) PhotoIds() []FileId {
	ret := make([]FileId, len(u.Message.Photo))
	for i, photo := range u.Message.Photo {
		ret[i] = FileId(photo.FileID)
	}
	return ret
}
