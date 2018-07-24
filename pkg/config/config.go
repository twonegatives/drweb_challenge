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
		cfg := viper.New()
		cfg.SetDefault("LISTEN", defaults.Listen)
		cfg.SetDefault("WRITE_TIMEOUT", defaults.WriteTimeout)
		cfg.SetDefault("READ_TIMEOUT", defaults.ReadTimeout)
		cfg.SetDefault("MAX_FILE_SIZE", defaults.MaxFileSize)
		cfg.AutomaticEnv()
	})

	return cfg
}

type configDefaults struct {
	Listen       string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	MaxFileSize  int
}

func getDefaults() *configDefaults {
	return &configDefaults{
		Listen:       ":80",
		WriteTimeout: 15,
		ReadTimeout:  15,
		MaxFileSize:  3 * 1000 * 1000,
	}
}
