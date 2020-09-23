package main

import (
	"context"
	"flag"
	"testing"

	"balance/internal/container"
	"balance/internal/exchangerates"
	"balance/internal/postgres"
	"balance/internal/services/core"
)

func TestHealth(t *testing.T) {
	flag.Parse()

	coreService := build(*debug)

	h := coreService.Health()

	if !h.Redis {
		t.Errorf("Healthcheck for Redis failed, errors: [%+v]", h.Errors)
	}

	if !h.Postgres {
		t.Errorf("Healthcheck for Postgres failed, errors: [%+v]", h.Errors)
	}
}

func build(debug bool) *core.Service {
	ctx := context.Background()

	test := true
	deps := container.Build(ctx, &debug, &test)

	err := postgres.CreateDatabaseSchema(deps.Postgres)
	if err != nil {
		deps.Logger.WithContext(ctx).Fatal(err)
	}

	exchangeRatesClient := exchangerates.NewClient(deps.HttpClient)
	exchangeRatesService := exchangerates.NewService(exchangeRatesClient, deps.Redis)

	coreService := core.New(ctx, deps.Redis, deps.Postgres, deps.Logger, deps.Stats, exchangeRatesService)

	return coreService
}
