package postgres

import (
	"context"

	"github.com/go-pg/pg/v10"
	"github.com/pkg/errors"

	"balance/internal/config"
)

func Client(ctx context.Context, cfg *config.Config) (*pg.DB, error) {
	opt, err := pg.ParseURL(cfg.PostgresUrl)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to Parse URL")
	}

	db := pg.Connect(opt).WithContext(ctx)

	_, err = db.Exec("select now()")
	if err != nil {
		return nil, errors.Wrap(err, "Unable to execute ping query")
	}

	return db, nil
}
