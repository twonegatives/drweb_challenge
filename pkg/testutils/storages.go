package testutils

import (
	"github.com/golang/mock/gomock"
	"github.com/twonegatives/drweb_challenge/pkg/drweb"
	"github.com/twonegatives/drweb_challenge/pkg/mocks"
	"github.com/twonegatives/drweb_challenge/pkg/storages"
)

func GenerateStorage(filename string, path string, err error, ctrl *gomock.Controller) drweb.Storage {
	pathgen := mocks.NewMockFilePathGenerator(ctrl)
	pathgen.EXPECT().Generate(filename).Return(path, err)
	return &storages.FileSystemStorage{
		FilePathGenerator: pathgen,
	}
}
