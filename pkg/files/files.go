package files

import (
	"io"
	"net/textproto"

	"github.com/pkg/errors"
	"github.com/twonegatives/drweb_challenge/pkg/drweb"
)

func NewFile(reader io.Reader, mimetype textproto.MIMEHeader, nameGenerator drweb.FileNameGenerator) (*drweb.File, error) {
	pipeReader, pipeWriter := io.Pipe()
	filenameReader := io.TeeReader(reader, pipeWriter)

	filename, err := nameGenerator.Generate(filenameReader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to generate filename")
	}

	file := drweb.File{
		Body:     pipeReader,
		MimeType: mimetype,
		Filename: filename,
	}

	return &file, nil
}
