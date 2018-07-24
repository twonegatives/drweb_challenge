package config

import (
	"time"

	"github.com/spf13/viper"
)

type configDefaults struct {
	Listen       string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

func getDefaults() *configDefaults {
	return &configDefaults{
		Listen:       ":80",
		WriteTimeout: 15,
		ReadTimeout:  15,
	}
}

func NewConfig() *viper.Viper {
	defaults := getDefaults()
	cfg := viper.New()
	cfg.SetDefault("LISTEN", defaults.Listen)
	cfg.SetDefault("WRITE_TIMEOUT", defaults.WriteTimeout)
	cfg.SetDefault("READ_TIMEOUT", defaults.ReadTimeout)
	cfg.AutomaticEnv()
	return cfg
}
