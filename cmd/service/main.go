package main

import (
	"context"
	"flag"

	"balance/internal/application"
)

var (
	debug = flag.Bool("debug", false, "run service in debug mode (with .env.debug)")
	test  = flag.Bool("test", false, "run service in test mode (with .env.test)")
)

func main() {
	flag.Parse()

	ctx := context.Background()

	app := application.Engine(ctx, *debug, *test)

	err := app.Router.Run()
	if err != nil {
		app.Container.Logger.Fatal(err)
	}
}
