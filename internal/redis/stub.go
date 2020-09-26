package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Stub struct {
	Cached map[string]string
}

func NewStub() *Stub {
	return &Stub{
		Cached: map[string]string{},
	}
}

func (s Stub) Get(_ context.Context, key string) (string, error) {
	result, ok := s.Cached[key]
	if !ok || result == "" {
		return "", redis.Nil
	}
	return result, nil
}

func (s Stub) Set(_ context.Context, key string, value interface{}, _ time.Duration) error {
	switch value.(type) {
	case string, []uint8:
		s.Cached[key] = fmt.Sprintf("%s", value)
		return nil
	default:
		marshaled, err := json.Marshal(value)
		if err != nil {
			return err
		}
		s.Cached[key] = string(marshaled)
		return nil
	}
}
