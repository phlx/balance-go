package container

import (
	"context"
	"net/http"
	"time"

	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"gopkg.in/alexcesaro/statsd.v2"

	"balance/internal/config"
	"balance/internal/logger"
	"balance/internal/postgres"
	"balance/internal/redis"
)

type Container struct {
	Logger     *logrus.Logger
	Stats      *statsd.Client
	Postgres   *pg.DB
	Redis      redis.Client
	HttpClient *http.Client
}

func Build(ctx context.Context, debug *bool, test *bool) *Container {
	log := logger.New(*debug)

	cfg, err := config.Load(*debug, *test)
	if err != nil {
		log.WithContext(ctx).Fatal(err)
	}

	if *test && cfg.Environment != "test" {
		err := errors.Errorf("aborted: test run in non-test environment (ENV=%s)", cfg.Environment)
		log.WithContext(ctx).Fatal(err)
	}

	stats, err := statsd.New(statsd.Address(cfg.StatsDAddress))
	if err != nil {
		log.WithContext(ctx).Fatal(err)
	}

	rds, err := redis.New(ctx, cfg)
	if err != nil {
		log.WithContext(ctx).Fatal(err)
	}

	db, err := postgres.Client(ctx, cfg)
	if err != nil {
		log.WithContext(ctx).Fatal(err)
	}
	db.AddQueryHook(postgres.NewHook(log))

	httpClient := http.Client{
		Transport: &http.Transport{
			Proxy:           nil,
			MaxIdleConns:    10,
			IdleConnTimeout: 30 * time.Second,
		},
		Timeout: 30 * time.Second,
	}

	return &Container{
		Logger:     log,
		Stats:      stats,
		Postgres:   db,
		Redis:      rds,
		HttpClient: &httpClient,
	}
}
