package configs

import (
	"net/url"
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Dsn      string
	Levellog zerolog.Level
}

func New() *Config {
	return &Config{
		Dsn:      getDsn(),
		Levellog: getLevelLog(),
	}
}

func getDsn() string {
	dsn := os.Getenv("postgres_url")
	_, err := url.Parse(dsn)
	if err != nil {
		log.Fatal().Str("dsn", dsn).Msg("not valid dsn")
	}

	return dsn
}

func getLevelLog() zerolog.Level {
	llevel := os.Getenv("log_level")
	switch lvl := strings.ToLower(llevel); lvl {
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "":
		log.Fatal().Str("log_level", llevel).Msg("level hasn't to be empty")
	default:
		log.Fatal().Str("log_level", llevel).Msg("level is not present in app")
	}

	return zerolog.Disabled
}
