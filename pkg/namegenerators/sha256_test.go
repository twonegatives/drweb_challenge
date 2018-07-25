package namegenerators_test

import (
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/twonegatives/drweb_challenge/pkg/namegenerators"
)

func TestGenerateWrongMimeType(t *testing.T) {
	generator := namegenerators.SHA256{}
	mimetype := "image/unexistant"
	reader := strings.NewReader("Some testing string")
	_, err := generator.Generate(reader, mimetype)

	assert.NotNil(t, err)
	assert.Equal(t, "could not find appropriate file extension", err.Error())
}

func TestGenerateReaderErrored(t *testing.T) {
	generator := namegenerators.SHA256{}
	mimetype := "image/png"
	reader := &errReader{}
	_, err := generator.Generate(reader, mimetype)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to hashify input stream")
}

func TestGenerateSuccess(t *testing.T) {
	generator := namegenerators.SHA256{}
	mimetype := "image/png"
	reader := strings.NewReader("Some testing string")
	hashstring, err := generator.Generate(reader, mimetype)

	assert.Nil(t, err)
	assert.Equal(t, "4859309121b35604ae3a848ac3a275b8d71410a1c09d9585c19ecea9fb84a2e2.png", hashstring)
}

type errReader struct {
}

func (e *errReader) Read(p []byte) (n int, err error) {
	return -1, errors.New("encountered error")
}
