package config

import (
	"flag"
	"testing"
)

var (
	// Fix of problem "flag provided but not defined: -debug" when running tests with "-args -debug"
	debug = flag.Bool("debug", false, "run service in debug mode (with .env.debug)")
)

func TestLoad(t *testing.T) {
	config, err := Load(*debug, true)
	if err != nil {
		t.Fatalf("Unable to load config: %s", err.Error())
	}
	if config.Environment == "" {
		t.Fatalf("Config environment is empty")
	}
}
