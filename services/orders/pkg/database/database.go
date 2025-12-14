package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/viper"
)

func NewClient(ctx context.Context) (*pgxpool.Pool, error) {
	host := viper.GetString("DB_ORDER_HOST")
	port := viper.GetString("DB_ORDER_PORT")
	user := viper.GetString("DB_ORDER_USERNAME")
	password := viper.GetString("DB_ORDER_PASSWORD")
	dbname := viper.GetString("DB_ORDER_NAME")

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		user, password, host, port, dbname,
	)

	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pgx config: %w", err)
	}

	maxConns := int32(viper.GetInt("DB_MAX_CONN_ORDER"))
	if maxConns == 0 {
		maxConns = 20
	}
	cfg.MaxConns = maxConns

	minConns := int32(viper.GetInt("DB_MIN_CONN_ORDER"))
	if minConns == 0 {
		minConns = 5
	}
	cfg.MinConns = minConns

	idleTime := viper.GetDuration("DB_MAX_IDLE_TIME")
	if idleTime == 0 {
		idleTime = 30 * time.Minute
	}
	cfg.MaxConnIdleTime = idleTime

	lifetime := viper.GetDuration("DB_MAX_LIFETIME")
	if lifetime == 0 {
		lifetime = 2 * time.Hour
	}
	cfg.MaxConnLifetime = lifetime

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create pgx pool: %w", err)
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return pool, nil
}
