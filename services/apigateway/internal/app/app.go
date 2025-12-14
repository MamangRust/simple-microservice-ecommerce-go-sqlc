package apps

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	_ "github.com/MamangRust/simple_microservice_ecommerce/apigateway/docs"
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/handler"
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/middlewares"
	mencache "github.com/MamangRust/simple_microservice_ecommerce/apigateway/internal/redis"
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/pkg/dotenv"
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/pkg/logger"
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/pkg/otel"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ServiceAddresses struct {
	Auth    string
	Role    string
	User    string
	Order   string
	Product string
}

func loadServiceAddresses() *ServiceAddresses {
	return &ServiceAddresses{
		Auth:    getEnvOrDefault("GRPC_AUTH_ADDR", "localhost:50051"),
		Role:    getEnvOrDefault("GRPC_ROLE_ADDR", "localhost:50052"),
		User:    getEnvOrDefault("GRPC_USER_ADDR", "localhost:50053"),
		Product: getEnvOrDefault("GRPC_PRODUCT_ADDR", "localhost:50054"),
		Order:   getEnvOrDefault("GRPC_ORDER_ADDR", "localhost:50055"),
	}
}

func createServiceConnections(addresses *ServiceAddresses, logger logger.LoggerInterface) (handler.ServiceConnections, error) {
	var connections handler.ServiceConnections

	conns := map[string]*string{
		"Auth":    &addresses.Auth,
		"Role":    &addresses.Role,
		"User":    &addresses.User,
		"Order":   &addresses.Order,
		"Product": &addresses.Product,
	}

	for name, addr := range conns {
		conn, err := createConnection(*addr, name, logger)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to connect to %s service, continuing with nil connection", name), zap.Error(err))
			conn = nil
		}
		switch name {
		case "Auth":
			connections.Auth = conn
		case "Role":
			connections.Role = conn
		case "User":
			connections.User = conn
		case "Order":
			connections.Order = conn
		case "Product":
			connections.Product = conn
		}
	}

	return connections, nil
}

// @title Ecommerce gRPC
// @version 1.0
// @description gRPC based Ecommerce service

// @host localhost:5000
// @BasePath /api/

// @securityDefinitions.apikey BearerAuth
// @in Header
// @name Authorization
type Client struct {
	App    *echo.Echo
	Logger logger.LoggerInterface
}

func (c *Client) Shutdown(ctx context.Context) error {
	c.Logger.Info("Shutting down API Gateway...")
	return c.App.Shutdown(ctx)
}

func RunClient() (*Client, func(), error) {
	flag.Parse()

	telemetry := otel.NewTelemetry(otel.Config{
		ServiceName:            "apigateway",
		ServiceVersion:         "1.0.0",
		Environment:            "production",
		Endpoint:               "otel-collector:4317",
		Insecure:               true,
		EnableRuntimeMetrics:   true,
		RuntimeMetricsInterval: 15 * time.Second,
	})

	if err := telemetry.Init(context.Background()); err != nil {
		return nil, nil, fmt.Errorf("failed to initialize telemetry: %w", err)
	}

	addresses := loadServiceAddresses()

	logger, err := logger.NewLogger("apigateway", telemetry.GetLogger())

	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	logger.Info("[apigateway] Loading environment variables...")
	if err := dotenv.Viper(); err != nil {
		return nil, nil, fmt.Errorf("failed to load .env file: %w", err)
	}

	logger.Info("[apigateway] Creating gRPC connections...")
	conns, err := createServiceConnections(addresses, logger)
	if err != nil {
		logger.Error("[apigateway] Failed to create some service connections, continuing with available connections", zap.Error(err))
	}

	myredis := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", viper.GetString("REDIS_HOST_APIGATEWAY"), viper.GetString("REDIS_PORT_APIGATEWAY")),
		Password:     viper.GetString("REDIS_PASSWORD_APIGATEWAY"),
		DB:           viper.GetInt("REDIS_DB_APIGATEWAY"),
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 3,
	})

	if err := myredis.Ping(context.Background()).Err(); err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
		return nil, nil, err
	}

	mencache := mencache.NewMencache(&mencache.Deps{
		Redis:  myredis,
		Logger: logger,
	})

	e := setupEcho(myredis, logger)

	depsHandler := &handler.Deps{
		E:                  e,
		ServiceConnections: &conns,
		Logger:             logger,
		Mencache:           mencache,
	}

	handler.NewHandler(depsHandler)

	errorHandler := middlewares.NewErrorHandlerMiddleware(logger)
	e.HTTPErrorHandler = errorHandler.ErrorHandler

	var wg sync.WaitGroup

	wg.Go(func() {
		logger.Info("[apigateway] Starting API Gateway on :5000")
		if err := e.Start(":5000"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Info("[apigateway] Echo server error:", zap.Error(err))
		}
	})

	shutdown := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		logger.Info("[apigateway] Shutting down API Gateway...")
		if err := e.Shutdown(ctx); err != nil {
			logger.Info("[apigateway] Echo shutdown failed:", zap.Error(err))
		}

		if err := telemetry.Shutdown(context.Background()); err != nil {
			logger.Error("Failed to shutdown tracer", zap.Error(err))
		}

		closeConnections(conns, logger)
		logger.Info("[apigateway] All connections closed.")
	}

	return &Client{App: e, Logger: logger}, shutdown, nil
}

func setupEcho(redis *redis.Client, logger logger.LoggerInterface) *echo.Echo {
	e := echo.New()

	// limiter := middlewares.NewRateLimiterMiddleware(redis, logger, 100, 150)
	e.Use(middleware.Recover(), middleware.Logger())

	e.Use(middleware.Recover(), middleware.Logger())

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:1420", "http://localhost:33451"},
		AllowMethods:     []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "X-API-Key"},
		AllowCredentials: true,
	}))

	middlewares.WebSecurityConfig(e)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	return e
}

func createConnection(address, serviceName string, logger logger.LoggerInterface) (*grpc.ClientConn, error) {
	logger.Info(fmt.Sprintf("Connecting to %s service at %s", serviceName, address))

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to connect to %s service", serviceName), zap.Error(err))
		return nil, err
	}
	return conn, nil
}

func closeConnections(conns handler.ServiceConnections, logger logger.LoggerInterface) {
	connsMap := map[string]*grpc.ClientConn{
		"Auth":    conns.Auth,
		"Role":    conns.Role,
		"User":    conns.User,
		"Order":   conns.Order,
		"Product": conns.Product,
	}

	for name, conn := range connsMap {
		if conn != nil {
			if err := conn.Close(); err != nil {
				logger.Error(fmt.Sprintf("Failed to close %s connection", name), zap.Error(err))
			}
		}
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
