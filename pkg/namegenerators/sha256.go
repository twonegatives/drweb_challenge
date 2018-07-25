package namegenerators

import (
	"crypto/sha256"
	"fmt"
	"io"
	"mime"

	"github.com/pkg/errors"
)

type SHA256 struct {
}

func (s *SHA256) Generate(input io.Reader, mimeType string) (string, error) {
	hasher := sha256.New()

	extensions, err := mime.ExtensionsByType(mimeType)
	if err != nil || len(extensions) == 0 {
		return "", errors.Wrap(err, "could not find appropriate file extension")
	}

	if _, err := io.Copy(hasher, input); err != nil {
		return "", errors.Wrap(err, "failed to hashify input stream")
	}

	filename := fmt.Sprintf("%x%s", hasher.Sum(nil), extensions[0])
	return filename, nil
}
