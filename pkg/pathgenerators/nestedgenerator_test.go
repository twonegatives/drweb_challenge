package pathgenerators_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twonegatives/drweb_challenge/pkg/pathgenerators"
)

type errCase struct {
	BasePath      string
	Levels        int
	FoldersLength int
	Filename      string
	ExpectedErr   string
}

func TestGenerateFailure(t *testing.T) {
	var objects = map[string]errCase{
		"wrong level": {
			BasePath:      "../../tmp",
			Levels:        -1,
			FoldersLength: 2,
			Filename:      "somefilename.png",
			ExpectedErr:   "nested path level should be a positive integer (given '-1')",
		},
		"wrong folder length": {
			BasePath:      "../../tmp",
			Levels:        3,
			FoldersLength: 0,
			Filename:      "somefilename.png",
			ExpectedErr:   "folder name should be at least 1 char long (given '0')",
		},
		"too short filename": {
			BasePath:      "../../tmp",
			Levels:        3,
			FoldersLength: 3,
			Filename:      "shrt.png",
			ExpectedErr:   "filename 'shrt.png' can't be used for path with nested structure of '3' levels and folder name size '3'",
		},
	}

	for testName, testObject := range objects {
		t.Run(testName, func(t *testing.T) {
			generator := pathgenerators.NestedGenerator{
				BasePath:     testObject.BasePath,
				Levels:       testObject.Levels,
				FolderLength: testObject.FoldersLength,
			}

			_, err := generator.Generate(testObject.Filename)
			assert.NotNil(t, err)
			assert.Equal(t, testObject.ExpectedErr, err.Error())
		})
	}
}

func TestGenerateSuccess(t *testing.T) {
	generator := pathgenerators.NestedGenerator{
		BasePath:     "../../tmp",
		Levels:       2,
		FolderLength: 2,
	}

	path, err := generator.Generate("somefilename.png")
	assert.Nil(t, err)
	assert.Equal(t, "../../tmp/so/me/somefilename.png", path)
}
