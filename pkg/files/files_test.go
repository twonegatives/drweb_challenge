package files_test

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/twonegatives/drweb_challenge/pkg/files"
	"github.com/twonegatives/drweb_challenge/pkg/mocks"
)

func TestNewFile(t *testing.T) {
	mimetype := "image/png"
	mockCtrl := gomock.NewController(t)
	reader := ioutil.NopCloser(bytes.NewReader([]byte("Some testing string")))
	defer reader.Close()
	nameGenerator := mocks.NewMockFileNameGenerator(mockCtrl)

	file := files.NewFile(reader, mimetype, nameGenerator)
	assert.Equal(t, mimetype, file.MimeType)
	assert.Equal(t, reader, file.Body)
	assert.Equal(t, nameGenerator, file.NameGenerator)
}
