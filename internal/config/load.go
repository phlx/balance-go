package config

import (
	"log"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"

	"balance/internal/pkg/filesystem"
)

func Load(debug bool, test bool) (*Config, error) {
	var err error
	configsDirectory, err := getConfigsDirectory()
	if err != nil {
		log.Fatal(err)
	}

	envFile := configsDirectory + getEnvFile(debug, test)

	cfg := Config{}

	err = godotenv.Load(envFile)
	if err != nil {
		return nil, err
	}

	err = env.Parse(&cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, err
}

func getConfigsDirectory() (string, error) {
	root, err := filesystem.GetRootDirectory()
	if err != nil {
		return "", err
	}

	sep := string(os.PathSeparator)

	return root + "configs" + sep, nil
}

func getEnvFile(debug, test bool) string {
	e := ".env"
	if debug {
		e = ".env.debug"
	}
	if test {
		e = ".env.test"
	}
	if debug && test {
		e = ".env.test.debug"
	}
	return e
}
