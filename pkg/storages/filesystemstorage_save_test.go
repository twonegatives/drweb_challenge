package storages_test

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/twonegatives/drweb_challenge/pkg/drweb"
	"github.com/twonegatives/drweb_challenge/pkg/mocks"
	"github.com/twonegatives/drweb_challenge/pkg/storages"
)

type staticFileNameGenerator struct {
	Name string
}

// NOTE: used to read file contents before generation of hashed name first
// TeeReader in storage#Save forces us to do it in order to write contents to file
func (g *staticFileNameGenerator) Generate(input io.Reader, extension string) (string, error) {
	ioutil.ReadAll(input)
	return g.Name, nil
}

func TestSaveFailure(t *testing.T) {
	t.Run("blank name generator", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		pathgen := mocks.NewMockFilePathGenerator(mockCtrl)
		pathgen.EXPECT().Generate(gomock.Any()).Times(0)

		storage := storages.FileSystemStorage{
			FilePathGenerator: pathgen,
		}

		_, err := storage.Save(&drweb.File{})
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "failed to save file without name generator")
	})

	t.Run("broken filename", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		namegen := mocks.NewMockFileNameGenerator(mockCtrl)
		namegen.EXPECT().Generate(gomock.Any(), gomock.Any()).Return("", errors.New("name generation error"))

		pathgen := mocks.NewMockFilePathGenerator(mockCtrl)
		pathgen.EXPECT().Generate(gomock.Any()).Times(0)

		storage := storages.FileSystemStorage{
			FilePathGenerator: pathgen,
		}

		_, err := storage.Save(&drweb.File{NameGenerator: namegen})
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "failed to generate filename")
	})

	t.Run("broken filepath", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		defer mockCtrl.Finish()

		namegen := mocks.NewMockFileNameGenerator(mockCtrl)
		namegen.EXPECT().Generate(gomock.Any(), gomock.Any()).Return("encrypted", nil)

		pathgen := mocks.NewMockFilePathGenerator(mockCtrl)
		pathgen.EXPECT().Generate("encrypted").Return("", errors.New("path generation error"))

		storage := storages.FileSystemStorage{
			FilePathGenerator: pathgen,
		}

		_, err := storage.Save(&drweb.File{NameGenerator: namegen})
		assert.NotNil(t, err)
		assert.Contains(t, err.Error(), "failed to generate filepath")
	})
}

func TestSaveSuccess(t *testing.T) {
	filename := "encrypted1"
	path := path.Join("../../tmp", filename)
	contents := []byte("File contents")
	namegen := &staticFileNameGenerator{Name: filename}

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	pathgen := mocks.NewMockFilePathGenerator(mockCtrl)
	pathgen.EXPECT().Generate(filename).Return(path, nil)

	storage := storages.FileSystemStorage{
		FileMode:          0700,
		FilePathGenerator: pathgen,
	}

	file := drweb.File{
		Body:          ioutil.NopCloser(bytes.NewReader(contents)),
		NameGenerator: namegen,
	}

	savedFileName, err := storage.Save(&file)
	defer os.Remove(path)

	if err != nil {
		t.Fatal(err)
	}

	bytes, err := ioutil.ReadFile(path)
	assert.Nil(t, err)
	assert.Equal(t, savedFileName, filename)
	assert.Equal(t, contents, bytes)
}
