package namegenerators

import (
	"fmt"
	"time"

	"github.com/twonegatives/drweb_challenge/pkg/drweb"
)

type LeadingSHA256WithUnixTime struct {
	LeadingSize int
}

func (l *LeadingSHA256WithUnixTime) Generate(file *drweb.File) (string, error) {
	encoder := SHA256Encoder{}
	leadingChars := make([]byte, l.LeadingSize)
	if _, err := file.Body.Read(leadingChars); err != nil {
		return file.Filename, err
	}

	filename := fmt.Sprintf("%x-%d", encoder.Encode(leadingChars), time.Now().UnixNano())
	return filename, nil
}
