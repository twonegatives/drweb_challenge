package files_test

import (
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/twonegatives/drweb_challenge/pkg/files"
	"github.com/twonegatives/drweb_challenge/pkg/mocks"
)

func TestNewFile(t *testing.T) {
	mimetype := "image/png"
	mockCtrl := gomock.NewController(t)
	reader := strings.NewReader("Some testing string")
	nameGenerator := mocks.NewMockFileNameGenerator(mockCtrl)

	file, err := files.NewFile(reader, mimetype, nameGenerator)
	assert.Nil(t, err)
	assert.Equal(t, mimetype, file.MimeType)
	assert.Equal(t, reader, file.Body)
	assert.Equal(t, nameGenerator, file.NameGenerator)
}
