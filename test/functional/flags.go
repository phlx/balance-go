package functional

import (
	"flag"
)

var (
	// Fix of error "flag provided but not defined: -debug" for all nested tests
	_ = flag.Bool("debug", false, "run service in debug mode (with .env.debug)")
)
