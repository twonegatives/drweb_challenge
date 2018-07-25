package files

import (
	"io"
	"net/textproto"

	"github.com/twonegatives/drweb_challenge/pkg/drweb"
)

func NewFile(reader io.Reader, mimetype textproto.MIMEHeader, nameGenerator drweb.FileNameGenerator) (*drweb.File, error) {
	file := drweb.File{
		Body:          reader,
		MimeType:      mimetype,
		NameGenerator: nameGenerator,
	}

	return &file, nil
}
