package files

import (
	"io"

	"github.com/twonegatives/drweb_challenge/pkg/drweb"
)

func NewFile(reader io.Reader, mimetype string, nameGenerator drweb.FileNameGenerator) (*drweb.File, error) {
	file := drweb.File{
		Body:          reader,
		MimeType:      mimetype,
		NameGenerator: nameGenerator,
	}

	return &file, nil
}
