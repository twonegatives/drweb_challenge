package file

import (
	"io/ioutil"
	"os"

	"github.com/twonegatives/drweb_challenge/pkg/drweb"
)

type File struct {
	Body     []byte
	Encoder  drweb.FileEncoder
	Uploader drweb.FileUploader
}

func (f *File) toHash() string {
	hashbytes := f.Encoder.Encode(f.Body)
	return string(hashbytes[:])
}

func (f *File) Save() (string, error) {
	filename := f.toHash()
	// TODO: call to Sync should be ensured
	return f.Uploader.Upload(f.Body, filename)
}

func LoadFile(hashstring string) (*File, error) {
	content, err := ioutil.ReadFile(hashstring)
	return &File{Body: content}, err
}

func DeleteFile(hashstring string) error {
	return os.Remove(hashstring)
}
