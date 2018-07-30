package namegenerators_test

import (
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/twonegatives/drweb_challenge/pkg/namegenerators"
)

type successCase struct {
	Content string
	Result  string
}

func TestGenerateSuccess(t *testing.T) {
	content := "Some testing string"
	result := "4859309121b35604ae3a848ac3a275b8d71410a1c09d9585c19ecea9fb84a2e2"
	generator := namegenerators.SHA256{}
	reader := strings.NewReader(content)
	hashstring, err := generator.Generate(reader)

	assert.Nil(t, err)
	assert.Equal(t, result, hashstring)
}

func TestGenerateReaderErrored(t *testing.T) {
	generator := namegenerators.SHA256{}
	reader := &errReader{}
	_, err := generator.Generate(reader)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to hashify input stream")
}

type errReader struct {
}

func (e *errReader) Read(p []byte) (n int, err error) {
	return -1, errors.New("encountered error")
}
