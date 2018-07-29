package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/twonegatives/drweb_challenge/pkg/config"
)

func TestGetConfig(t *testing.T) {
	cfg := config.GetConfig()
	assert.NotNil(t, cfg)
	assert.Equal(t, ":80", cfg.GetString("LISTEN"))
}
