package namegenerators_test

import (
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/twonegatives/drweb_challenge/pkg/namegenerators"
)

type successCase struct {
	Content   string
	Extension string
	Result    string
}

func TestGenerateSuccess(t *testing.T) {
	var objects = map[string]successCase{
		"with extension": {
			Content:   "Some testing string",
			Extension: ".png",
			Result:    "4859309121b35604ae3a848ac3a275b8d71410a1c09d9585c19ecea9fb84a2e2.png",
		},
		"without extension": {
			Content:   "Some testing string",
			Extension: "",
			Result:    "4859309121b35604ae3a848ac3a275b8d71410a1c09d9585c19ecea9fb84a2e2",
		},
	}

	for testName, testObject := range objects {
		t.Run(testName, func(t *testing.T) {
			generator := namegenerators.SHA256{}
			reader := strings.NewReader(testObject.Content)
			hashstring, err := generator.Generate(reader, testObject.Extension)

			assert.Nil(t, err)
			assert.Equal(t, testObject.Result, hashstring)
		})
	}
}

func TestGenerateReaderErrored(t *testing.T) {
	generator := namegenerators.SHA256{}
	extension := ".png"
	reader := &errReader{}
	_, err := generator.Generate(reader, extension)

	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "failed to hashify input stream")
}

type errReader struct {
}

func (e *errReader) Read(p []byte) (n int, err error) {
	return -1, errors.New("encountered error")
}
