package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"

	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/dotenv"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/logger"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/otel"
	"github.com/pressly/goose/v3"
	"github.com/spf13/viper"
)

const (
	dialect = "postgres"
)

var (
	flags = flag.NewFlagSet("migrate", flag.ExitOnError)
	dir   = flags.String("dir", "./pkg/database/migrations", "directory with migration files")
)

func main() {
	flags.Usage = usage
	if err := flags.Parse(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing flags: %v\n", err)
		os.Exit(1)
	}

	telemetry := otel.NewTelemetry(otel.Config{
		ServiceName:    "auth-migrate-service",
		ServiceVersion: "1.0.0",
		Environment:    "production",
		Endpoint:       "otel-collector:4317",
		Insecure:       true,
	})

	if err := telemetry.Init(context.Background()); err != nil {
		panic(fmt.Errorf("failed to initialize telemetry: %w", err))
	}

	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := telemetry.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "Error during telemetry shutdown: %v\n", err)
		}
	}()

	logger, err := logger.NewLogger("auth-migrate-service", telemetry.GetLogger())

	if err != nil {
		logger.Error("failed to initialize logger", zap.Error(err))
	}

	args := flags.Args()
	if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
		flags.Usage()
		return
	}

	command := args[0]

	err = dotenv.Viper()
	if err != nil {
		logger.Fatal("Error loading environment variables:", zap.Error(err))
	}

	connStr := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		viper.GetString("DB_AUTH_HOST"),
		viper.GetString("DB_AUTH_PORT"),
		viper.GetString("DB_AUTH_USERNAME"),
		viper.GetString("DB_AUTH_NAME"),
		viper.GetString("DB_AUTH_PASSWORD"),
	)

	db, err := goose.OpenDBWithDriver(dialect, connStr)
	if err != nil {
		logger.Fatal("Error opening database: ", zap.Error(err))
	}

	defer func() {
		if err := db.Close(); err != nil {
			logger.Fatal("Error closing database: ", zap.Error(err))
		}
	}()

	if err := goose.RunContext(context.Background(), command, db, *dir, args[1:]...); err != nil {
		logger.Fatal("Migration failed: ", zap.Error(err))
	}

	logger.Info("Migrate selesai")
}

func usage() {
	fmt.Println(usagePrefix)
	flags.PrintDefaults()
	fmt.Println(usageCommands)
}

var (
	usagePrefix = `Usage: migrate COMMAND
Examples:
    migrate status
`

	usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-by-one            Migrate the DB up by 1
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    reset                Roll back all migrations
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with the current timestamp
    fix                  Apply sequential ordering to migrations`
)
