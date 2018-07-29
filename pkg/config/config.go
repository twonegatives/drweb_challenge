package config

import (
	"sync"
	"time"

	"github.com/spf13/viper"
)

var cfg *viper.Viper
var once sync.Once

func GetConfig() *viper.Viper {
	once.Do(func() {
		defaults := getDefaults()
		cfg = viper.New()
		cfg.SetDefault("LISTEN", defaults.Listen)
		cfg.SetDefault("WRITE_TIMEOUT", defaults.WriteTimeout)
		cfg.SetDefault("READ_TIMEOUT", defaults.ReadTimeout)
		cfg.SetDefault("PATH_NESTED_LEVELS", defaults.PathNestedLevels)
		cfg.SetDefault("PATH_NESTED_FOLDERS_LENGTH", defaults.PathNestedFoldersLength)
		cfg.SetDefault("PATH_BASE", defaults.PathBase)
		cfg.SetDefault("STORAGE_FILE_MODE", defaults.StorageFileMode)
		cfg.AutomaticEnv()
	})

	return cfg
}

type configDefaults struct {
	Listen                  string
	ReadTimeout             time.Duration
	WriteTimeout            time.Duration
	PathNestedLevels        int
	PathNestedFoldersLength int
	PathBase                string
	StorageFileMode         int
}

func getDefaults() *configDefaults {
	return &configDefaults{
		Listen:       ":80",
		WriteTimeout: 15,
		ReadTimeout:  15,
		// NOTE: we use double folder nesting here in order to overcome
		// issue with too much files lying in a single folder.
		PathNestedLevels:        2,
		PathNestedFoldersLength: 2,
		PathBase:                ".",
		StorageFileMode:         0755,
	}
}
