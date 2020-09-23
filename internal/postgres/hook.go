package postgres

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"
)

type Hook struct {
	logger *logrus.Logger
}

func NewHook(logger *logrus.Logger) *Hook {
	return &Hook{
		logger: logger,
	}
}

func (h *Hook) BeforeQuery(ctx context.Context, e *pg.QueryEvent) (context.Context, error) {
	return ctx, nil
}

func (h *Hook) AfterQuery(_ context.Context, e *pg.QueryEvent) error {
	formatted, _ := e.FormattedQuery()
	h.logger.WithField("formatted", string(formatted)).Debug("[SQL]")
	return nil
}
