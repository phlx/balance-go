package core

import (
	"context"
	"crypto/md5" // #nosec
	"time"

	"github.com/go-pg/pg/v10"

	"balance/internal/redis"
)

type Result struct {
	Postgres bool
	Redis    bool
	Errors   []string
}

func (s *Service) Health() Result {
	result := Result{
		Postgres: true,
		Redis:    true,
		Errors:   []string{},
	}
	var err error
	err = checkPostgres(s.context, s.postgres)
	if err != nil {
		result.Postgres = false
		result.Errors = append(result.Errors, err.Error())
	}
	err = checkRedis(s.context, s.redis)
	if err != nil {
		result.Redis = false
		result.Errors = append(result.Errors, err.Error())
	}
	return result
}

func checkPostgres(ctx context.Context, p *pg.DB) error {
	_, err := p.ExecOneContext(ctx, "select now()")
	return err
}

func checkRedis(ctx context.Context, r redis.Client) error {
	t := time.Now().String()
	hash := md5.Sum([]byte(t)) // #nosec
	key := "temp:" + string(hash[:])
	err := r.Set(ctx, key, t, 1*time.Second)
	if err != nil {
		return err
	}
	_, err = r.Get(ctx, key)
	return err
}
