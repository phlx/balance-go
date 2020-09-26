package application

import (
	"context"

	"github.com/gin-gonic/gin"

	"balance/internal/container"
	"balance/internal/exchangerates"
	"balance/internal/handlers/balance"
	"balance/internal/handlers/give"
	"balance/internal/handlers/health"
	"balance/internal/handlers/move"
	"balance/internal/handlers/root"
	"balance/internal/handlers/take"
	"balance/internal/handlers/transactions"
	"balance/internal/middlewares"
	"balance/internal/postgres"
	"balance/internal/services/core"
	"balance/internal/validations"
)

func Engine(ctx context.Context, debug, test bool) *Application {
	deps := container.Build(ctx, debug, test)

	err := postgres.CreateDatabaseSchema(deps.Postgres)
	if err != nil {
		deps.Logger.WithContext(ctx).Fatal(err)
	}

	if !debug {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.New()
	router.Use(middlewares.Logger(deps.Logger), gin.Recovery())
	router.Use(middlewares.Cors())
	router.Use(middlewares.Idempotency(ctx, deps.Postgres, deps.Logger))

	validations.RegisterCurrency(deps.Logger)

	router.GET("/", root.Controller())

	exchangeRatesClient := exchangerates.NewClient(deps.HttpClient)
	exchangeRatesService := exchangerates.NewService(exchangeRatesClient, deps.Redis)

	coreService := core.New(ctx, deps.Redis, deps.Postgres, deps.Logger, deps.Metrics, exchangeRatesService)
	router.GET("/_health", health.Controller(coreService, deps.Metrics))
	router.GET("/v1/balance", balance.Controller(coreService, deps.Metrics))
	router.POST("/v1/give", give.Controller(coreService, deps.Metrics))
	router.POST("/v1/take", take.Controller(coreService, deps.Metrics))
	router.POST("/v1/move", move.Controller(coreService, deps.Metrics))
	router.GET("/v1/transactions", transactions.Controller(coreService, deps.Metrics))

	return &Application{
		Container:            deps,
		Router:               router,
		CoreService:          coreService,
		ExchangeRatesService: exchangeRatesService,
	}
}
