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
	hasher := sha256.New()

	if _, err := io.Copy(hasher, input); err != nil {
		return "", errors.Wrap(err, "failed to hashify input stream")
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}
