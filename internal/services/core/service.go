package core

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"

	"balance/internal/exchangerates"
	"balance/internal/metrics"
	"balance/internal/redis"
)

type Service struct {
	context              context.Context
	redis                redis.Client
	stats                metrics.Client
	postgres             *pg.DB
	logger               *logrus.Logger
	exchangeRatesService *exchangerates.Service
}

func New(
	context context.Context,
	redis redis.Client,
	postgres *pg.DB,
	logger *logrus.Logger,
	stats metrics.Client,
	exchangeRatesService *exchangerates.Service,
) *Service {
	return &Service{
		context:              context,
		redis:                redis,
		postgres:             postgres,
		logger:               logger,
		stats:                stats,
		exchangeRatesService: exchangeRatesService,
	}
}
