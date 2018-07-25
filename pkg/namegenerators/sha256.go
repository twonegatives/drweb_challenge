package namegenerators

import (
	"crypto/sha256"
	"fmt"
	"io"

	"github.com/pkg/errors"
)

type SHA256 struct {
}

func (s *SHA256) Generate(input io.Reader) (string, error) {
	h := sha256.New()
	if _, err := io.Copy(h, input); err != nil {
		return "", errors.Wrap(err, "failed to hashify input stream")
	}

	filename := fmt.Sprintf("%x", h.Sum(nil))
	return filename, nil
}
