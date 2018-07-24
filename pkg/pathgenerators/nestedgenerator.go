package pathgenerators

import (
	"fmt"
	"path"
)

type NestedGenerator struct {
	BasePath     string
	Levels       int
	FolderLength int
}

func (g *NestedGenerator) Generate(filename string) (string, error) {
	resultPath := g.BasePath

	if g.Levels < 0 {
		return resultPath, fmt.Errorf("nested path level should be a positive integer (given '%d')", g.Levels)
	}

	if g.Levels > 0 && g.FolderLength < 1 {
		return resultPath, fmt.Errorf("folder name should be at least 1 char long (given '%d')", g.FolderLength)
	}

	if len(filename) < g.Levels*g.FolderLength {
		errString := "filename '%s' can't be used for path with nested structure of '%d' levels and folder name size '%d'"
		return resultPath, fmt.Errorf(errString, filename, g.Levels, g.FolderLength)
	}

	for i := 0; i < g.Levels; i++ {
		lower := i * g.FolderLength
		upper := lower + g.FolderLength
		resultPath = path.Join(resultPath, filename[lower:upper])
	}

	return resultPath, nil
}
