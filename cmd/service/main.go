package main

import (
	"context"
	"flag"

	"github.com/gin-gonic/gin"

	"balance/internal/container"
	"balance/internal/controllers/balance"
	"balance/internal/controllers/give"
	"balance/internal/controllers/health"
	"balance/internal/controllers/move"
	"balance/internal/controllers/root"
	"balance/internal/controllers/take"
	"balance/internal/controllers/transactions"
	"balance/internal/exchangerates"
	"balance/internal/middlewares"
	"balance/internal/postgres"
	"balance/internal/services/core"
	"balance/internal/validations"
)

var (
	debug = flag.Bool("debug", false, "run service in debug mode (with .env.debug)")
	test  = flag.Bool("test", false, "run service in test mode (with .env.test)")
)

func main() {
	flag.Parse()

	ctx := context.Background()

	deps := container.Build(ctx, debug, test)

	err := postgres.CreateDatabaseSchema(deps.Postgres)
	if err != nil {
		deps.Logger.WithContext(ctx).Fatal(err)
	}

	if !*debug {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(middlewares.Logger(deps.Logger), gin.Recovery())
	router.Use(middlewares.Cors())
	router.Use(middlewares.Idempotency(ctx, deps.Redis, deps.Logger))

	validations.RegisterCurrency(deps.Logger)

	router.GET("/", root.Controller())

	exchangeRatesClient := exchangerates.NewClient(deps.HttpClient)
	exchangeRatesService := exchangerates.NewService(exchangeRatesClient, deps.Redis)

	coreService := core.New(ctx, deps.Redis, deps.Postgres, deps.Logger, deps.Stats, exchangeRatesService)
	router.GET("/_health", health.Controller(coreService, deps.Stats))
	router.GET("/v1/balance", balance.Controller(coreService, deps.Stats))
	router.POST("/v1/give", give.Controller(coreService, deps.Stats))
	router.POST("/v1/take", take.Controller(coreService, deps.Stats))
	router.POST("/v1/move", move.Controller(coreService, deps.Stats))
	router.GET("/v1/transactions", transactions.Controller(coreService, deps.Stats))

	err = router.Run()
	if err != nil {
		deps.Logger.Fatal(err)
	}
}
