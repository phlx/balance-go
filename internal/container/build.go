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
	"balance/internal/metrics"
	"balance/internal/postgres"
	"balance/internal/redis"
)

type Container struct {
	Logger     *logrus.Logger
	Postgres   *pg.DB
	Metrics    metrics.Client
	Redis      redis.Client
	HttpClient *http.Client
}

func Build(ctx context.Context, debug bool, test bool) *Container {
	loggerClient := logger.New(debug)

	cfg, err := config.Load(debug, test)
	if err != nil {
		loggerClient.WithContext(ctx).Fatal(errors.Wrap(err, "unable to load config"))
	}

	if test && cfg.Environment != "test" {
		err := errors.Errorf("aborted: test run in non-test environment (ENV=%s)", cfg.Environment)
		loggerClient.WithContext(ctx).Fatal(err)
	}

	var metricsClient metrics.Client
	if test {
		metricsClient = metrics.NewStub()
	} else {
		stats, err := statsd.New(statsd.Address(cfg.StatsDAddress))
		if err != nil {
			err = errors.Wrapf(err, "unable to create StatsD client with addr %s", cfg.StatsDAddress)
			loggerClient.WithContext(ctx).Fatal(err)
		}
		metricsClient = metrics.New(stats)
	}

	var redisClient redis.Client
	if test {
		redisClient = redis.NewStub()
	} else {
		redisClient, err = redis.New(ctx, cfg.RedisAddress)
		if err != nil {
			err = errors.Wrapf(err, "unable to create Redis client with addr %s", cfg.RedisAddress)
			loggerClient.WithContext(ctx).Fatal(err)
		}
	}

	postgresClient, err := postgres.Client(ctx, cfg.PostgresUrl)
	if err != nil {
		err = errors.Wrapf(err, "unable to create Postgres client with URL %s", cfg.PostgresUrl)
		loggerClient.WithContext(ctx).Fatal(err)
	}
	postgresClient.AddQueryHook(postgres.NewHook(loggerClient))

	httpClient := &http.Client{
		Transport: &http.Transport{
			Proxy:           nil,
			MaxIdleConns:    10,
			IdleConnTimeout: 30 * time.Second,
		},
		Timeout: 30 * time.Second,
	}

	return &Container{
		Logger:     loggerClient,
		Metrics:    metricsClient,
		Postgres:   postgresClient,
		Redis:      redisClient,
		HttpClient: httpClient,
	}
}
