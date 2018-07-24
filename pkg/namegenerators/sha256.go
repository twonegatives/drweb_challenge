package namegenerators

import (
	"fmt"
	"io/ioutil"

	"github.com/twonegatives/drweb_challenge/pkg/drweb"
)

type SHA256 struct {
}

func (s *SHA256) Generate(file *drweb.File) (string, error) {
	encoder := SHA256Encoder{}
	contents, err := ioutil.ReadAll(file.Body)
	if err != nil {
		return file.Filename, err
	}

	filename := fmt.Sprintf("%x", encoder.Encode(contents))
	return filename, nil
}
