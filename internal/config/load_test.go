package config

import (
	"testing"
)

func TestLoad(t *testing.T) {
	config, err := Load(false, false)
	if err != nil {
		t.Fatalf("Unable to load config: %s", err.Error())
	}
	if config.Environment == "" {
		t.Fatalf("Config environment is empty")
	}
}
