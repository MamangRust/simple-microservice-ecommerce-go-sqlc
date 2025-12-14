package apps

import (
	"context"
	"flag"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/handler"
	mencache "github.com/MamangRust/simple_microservice_ecommerce/product/internal/redis"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"

	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/repository"
	"github.com/MamangRust/simple_microservice_ecommerce/product/internal/service"
	"github.com/MamangRust/simple_microservice_ecommerce/product/pkg/database"
	db "github.com/MamangRust/simple_microservice_ecommerce/product/pkg/database/schema"
	"github.com/MamangRust/simple_microservice_ecommerce/product/pkg/database/seeder"
	"github.com/MamangRust/simple_microservice_ecommerce/product/pkg/dotenv"
	"github.com/MamangRust/simple_microservice_ecommerce/product/pkg/logger"
	"github.com/MamangRust/simple_microservice_ecommerce/product/pkg/otel"
	pbproduct "github.com/MamangRust/simple_microservice_ecommerce_pb/product"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

var (
	port int
)

func init() {
	port = viper.GetInt("GRPC_PRODUCT_PORT")

	if port == 0 {
		port = 50054
	}

	flag.IntVar(&port, "port", port, "gRPC server port")
}

type Server struct {
	DB       *db.Queries
	Services service.Service
	Handlers handler.Handler
	Ctx      context.Context
	Logger   logger.LoggerInterface
}

func NewServer(ctx context.Context) (*Server, func(context.Context) error, error) {
	flag.Parse()

	telemetry := otel.NewTelemetry(otel.Config{
		ServiceName:            "product-service",
		ServiceVersion:         "1.0.0",
		Environment:            "production",
		Endpoint:               "otel-collector:4317",
		Insecure:               true,
		EnableRuntimeMetrics:   true,
		RuntimeMetricsInterval: 15 * time.Second,
	})

	if err := telemetry.Init(ctx); err != nil {
		return nil, nil, fmt.Errorf("failed to initialize telemetry: %w", err)
	}

	logger, err := logger.NewLogger("product-service", telemetry.GetLogger())

	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	if err := dotenv.Viper(); err != nil {
		logger.Fatal("Failed to load .env file", zap.Error(err))
	}
	logger.Info("Successfully loaded .env file")

	conn, err := database.NewClient(ctx)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	logger.Info("Successfully connected to the database")
	DB := db.New(conn)

	repositories := repository.NewRepositories(DB)
	logger.Info("Repository layer initialized")

	myredis := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", viper.GetString("REDIS_HOST_PRODUCT"), viper.GetString("REDIS_PORT_PRODUCT")),
		Password:     viper.GetString("REDIS_PASSWORD_PRODUCT"),
		DB:           viper.GetInt("REDIS_DB_PRODUCT"),
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolSize:     10,
		MinIdleConns: 3,
	})

	if err := myredis.Ping(ctx).Err(); err != nil {
		logger.Fatal("Failed to connect to Redis", zap.Error(err))
		return nil, nil, err
	}

	mencache := mencache.NewMencache(&mencache.Deps{
		Redis:  myredis,
		Logger: logger,
	})

	services := service.NewService(&service.Deps{
		Repositories: repositories,
		Mencache:     mencache,
		Logger:       logger,
	})

	logger.Info("Service layer initialized")

	handlers := handler.NewHandler(services, logger)
	logger.Info("Handler layer initialized")

	logger.Info("Server dependencies initialized successfully")
	return &Server{
		DB:       DB,
		Services: services,
		Logger:   logger,
		Handlers: handlers,
		Ctx:      ctx,
	}, telemetry.Shutdown, nil
}

func (s *Server) Run() {
	s.Logger.Info("Starting gRPC server...")

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		s.Logger.Fatal("Failed to start TCP listener", zap.Error(err), zap.Int("port", port))
	}
	s.Logger.Info("TCP listener started successfully", zap.String("address", lis.Addr().String()))

	grpcServer := grpc.NewServer()

	s.Logger.Info("Registering gRPC services...")
	pbproduct.RegisterProductQueryServiceServer(grpcServer, s.Handlers)
	pbproduct.RegisterProductCommandServiceServer(grpcServer, s.Handlers)
	s.Logger.Info("gRPC services registered successfully")

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.Logger.Info("gRPC server is starting to serve requests", zap.Int("port", port))
		if err := grpcServer.Serve(lis); err != nil {
			s.Logger.Fatal("Failed to start gRPC server", zap.Error(err))
		}
	}()

	wg.Wait()
}

func (s *Server) Seed() {
	dbSeederEnabled := viper.GetString("DB_SEEDER")
	s.Logger.Info("Checking database seeder flag", zap.String("DB_SEEDER", dbSeederEnabled))

	if dbSeederEnabled == "true" {
		s.Logger.Info("Database seeder is enabled, starting seeding process...")

		seeder := seeder.NewProductSeeder(s.DB, s.Logger)

		if err := seeder.SeedAll(s.Ctx); err != nil {
			s.Logger.Error("Failed to seed database", zap.Error(err))
			return
		}

		s.Logger.Info("Database seeding completed successfully.")
	} else {
		s.Logger.Info("Database seeder is disabled.")
	}
}
