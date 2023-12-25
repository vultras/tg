package tg

import (
	"bufio"
	"errors"
	"io"
	//"os"
	//"path/filepath"

	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type FileConfig = tgbotapi.FileConfig
type PhotoConfig = tgbotapi.PhotoConfig
type FileType int

const (
	NoFileType FileType = iota
	PhotoFileType
	DocumentFileType
)

var (
	UnknownFileTypeErr = errors.New("unknown file type")
)

// The type implements the structure to easily send
// files to the client.
type File struct {
	*MessageCompo
	name string
	reader io.Reader
	upload bool
	typ     FileType
	data, caption string
}

// Create the new file with the specified reader.
// By default it NeedsUpload is set to true.
func NewFile(reader io.Reader) *File {
	ret := &File{}

	ret.MessageCompo = NewMessage("")
	ret.reader = reader
	ret.upload = true

	return ret
}

func (f *File) Name(name string) *File {
	f.name = name
	return f
}

func (f *File) withType(typ FileType) *File {
	f.typ = typ
	return f
}

// Get the file type.
func (f *File) Type() FileType {
	return f.typ
}

// Set the file type to PhotoFileType.
func (f *File) Photo() *File {
	return f.withType(PhotoFileType)
}

func (f *File) Document() *File {
	return f.withType(DocumentFileType)
}

// Set the file caption.
func (f *File) Caption(caption string) *File {
	f.caption = caption
	return f
}

// Specifiy whether the file needs to be uploaded to Telegram.
func (f *File) Upload(upload bool) *File {
	f.upload = upload
	return f
}

// Set the data to return via SendData()
func (f *File) Data(data string) *File {
	f.data = data
	return f
}

func (f *File) NeedsUpload() bool {
	return f.upload
}

func (f *File) UploadData() (string, io.Reader, error) {
	// Bufferizing the reader
	// to make it faster.
	bufRd := bufio.NewReader(f.reader)
	fileName := f.name

	return fileName, bufRd, nil
}

func (f *File) SendData() string {
	return f.data
}

func (f *File) SendConfig(
	sid SessionId, bot *Bot,
) (*SendConfig) {
	var config SendConfig
	cid := sid.ToApi()

	switch f.Type() {
	case PhotoFileType:
		photo := tgbotapi.NewPhoto(cid, f)
		photo.Caption = f.caption

		config.Photo = &photo
	case DocumentFileType:
		doc := tgbotapi.NewDocument(sid.ToApi(), f)
		doc.Caption = f.caption
		config.Document = &doc
	default:
		panic(UnknownFileTypeErr)
	}


	return &config
}
