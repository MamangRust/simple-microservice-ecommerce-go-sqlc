package database

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/logger"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func NewClient(logger logger.LoggerInterface) (*sql.DB, error) {
	dbDriver := viper.GetString("DB_ORDER_DRIVER")
	if dbDriver == "" {
		dbDriver = viper.GetString("DB_DRIVER")
	}

	logger.Info("Initializing database client", zap.String("driver", dbDriver))

	var connStr string
	var logFields []zap.Field

	switch dbDriver {
	case "postgres":
		host := viper.GetString("DB_ORDER_HOST")
		port := viper.GetString("DB_ORDER_PORT")
		user := viper.GetString("DB_ORDER_USERNAME")
		dbname := viper.GetString("DB_ORDER_NAME")
		password := viper.GetString("DB_ORDER_PASSWORD")

		connStr = fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
			host, port, user, dbname, password,
		)
		logFields = []zap.Field{
			zap.String("host", host),
			zap.String("port", port),
			zap.String("dbname", dbname),
			zap.String("user", user),
		}

	case "mysql":
		host := viper.GetString("DB_ORDER_HOST")
		port := viper.GetString("DB_ORDER_PORT")
		user := viper.GetString("DB_ORDER_USERNAME")
		dbname := viper.GetString("DB_ORDER_NAME")
		password := viper.GetString("DB_ORDER_PASSWORD")

		connStr = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
			user, password, host, port, dbname,
		)
		logFields = []zap.Field{
			zap.String("host", host),
			zap.String("port", port),
			zap.String("dbname", dbname),
			zap.String("user", user),
		}
	default:
		logger.Fatal("Unsupported database driver", zap.String("driver", dbDriver))
		return nil, fmt.Errorf("unsupported database driver: %s", dbDriver)
	}

	logger.Debug("Attempting to connect to database", logFields...)

	db, err := sql.Open(dbDriver, connStr)
	if err != nil {
		logger.Fatal("Failed to open database connection", zap.Error(err), zap.String("driver", dbDriver))
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	maxOpenConns := viper.GetInt("DB_MAX_OPEN_CONNS")
	if maxOpenConns == 0 {
		maxOpenConns = 25
	}
	db.SetMaxOpenConns(maxOpenConns)

	maxIdleConns := viper.GetInt("DB_MAX_IDLE_CONNS")
	if maxIdleConns == 0 {
		maxIdleConns = 5
		logger.Debug("Using default value for DB_MAX_IDLE_CONNS", zap.Int("default", maxIdleConns))
	}
	db.SetMaxIdleConns(maxIdleConns)

	connMaxLifetime := viper.GetDuration("DB_CONN_MAX_LIFETIME")
	if connMaxLifetime == 0 {
		connMaxLifetime = time.Hour
		logger.Debug("Using default value for DB_CONN_MAX_LIFETIME", zap.Duration("default", connMaxLifetime))
	}
	db.SetConnMaxLifetime(connMaxLifetime)

	logger.Info("Database connection established successfully",
		zap.String("driver", dbDriver),
	)

	logger.Debug("Database connection pool configuration",
		zap.Int("max_open_conns", maxOpenConns),
		zap.Int("max_idle_conns", maxIdleConns),
		zap.Duration("conn_max_lifetime", connMaxLifetime),
	)

	return db, nil
}
