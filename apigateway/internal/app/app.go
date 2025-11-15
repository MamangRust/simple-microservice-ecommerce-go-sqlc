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
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/pkg/dotenv"
	"github.com/MamangRust/simple_microservice_ecommerce/apigateway/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
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

	addresses := loadServiceAddresses()

	logger, err := logger.NewLogger("apigateway")

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

	e := setupEcho()

	depsHandler := &handler.Deps{
		E:                  e,
		ServiceConnections: &conns,
		Logger:             logger,
	}

	handler.NewHandler(depsHandler)

	errorHandler := middlewares.NewErrorHandlerMiddleware(logger)
	e.HTTPErrorHandler = errorHandler.ErrorHandler

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		logger.Info("[apigateway] Starting API Gateway on :5000")
		if err := e.Start(":5000"); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Info("[apigateway] Echo server error:", zap.Error(err))
		}
	}()

	shutdown := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		logger.Info("[apigateway] Shutting down API Gateway...")
		if err := e.Shutdown(ctx); err != nil {
			logger.Info("[apigateway] Echo shutdown failed:", zap.Error(err))
		}

		closeConnections(conns, logger)
		logger.Info("[apigateway] All connections closed.")
	}

	return &Client{App: e, Logger: logger}, shutdown, nil
}

func setupEcho() *echo.Echo {
	e := echo.New()

	limiter := middlewares.NewRateLimiter(20, 50)
	e.Use(limiter.Limit, middleware.Recover(), middleware.Logger())

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
