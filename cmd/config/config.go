package config

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

const (
	XmlUrl = "XML_URL"
	DBUrl  = "DB_URL"
)

func LoadConfig() {
	viper.AutomaticEnv()
	viper.SetDefault("XML_URL", "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml")
	viper.SetDefault("DB_URL", "postgresql://postgres:postgres@localhost:5432/exchange_rates?sslmode=disable")
}

func SetupDatabase() (*sqlx.DB, error) {
	dsn := viper.GetString(DBUrl)

	if len(dsn) == 0 {
		return nil, fmt.Errorf("invalid db dsn: %s", dsn)
	}

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	return db, nil
}

func SetupLogger() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	slog.SetDefault(logger)

	return logger
}
