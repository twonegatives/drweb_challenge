package pathgenerators_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twonegatives/drweb_challenge/pkg/pathgenerators"
)

func TestGenerateWrongLevel(t *testing.T) {
	generator := pathgenerators.NestedGenerator{
		BasePath:     "../../tmp",
		Levels:       -1,
		FolderLength: 2,
	}

	_, err := generator.Generate("somefilename.png")
	assert.NotNil(t, err)
	assert.Equal(t, "nested path level should be a positive integer (given '-1')", err.Error())
}

func TestGenerateWrongFolderLength(t *testing.T) {
	generator := pathgenerators.NestedGenerator{
		BasePath:     "../../tmp",
		Levels:       3,
		FolderLength: 0,
	}

	_, err := generator.Generate("somefilename.png")
	assert.NotNil(t, err)
	assert.Equal(t, "folder name should be at least 1 char long (given '0')", err.Error())
}

func TestGenerateTooShortFilename(t *testing.T) {
	generator := pathgenerators.NestedGenerator{
		BasePath:     "../../tmp",
		Levels:       3,
		FolderLength: 3,
	}

	_, err := generator.Generate("shrt.png")
	assert.NotNil(t, err)
	assert.Equal(t, "filename 'shrt.png' can't be used for path with nested structure of '3' levels and folder name size '3'", err.Error())
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
