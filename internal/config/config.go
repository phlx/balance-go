package config

type Config struct {
	Environment   string `env:"ENV" envDefault:"development"`
	RedisAddress  string `env:"REDIS_ADDR" envDefault:"redis:6379"`
	PostgresUrl   string `env:"POSTGRES_URL" envDefault:"postgres://master:password@postgres:5432/app?sslmode=disable"`
	StatsDAddress string `env:"STATSD_ADDR" envDefault:"graphite:8125"`
}
