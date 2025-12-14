package apps

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	grpcclient "github.com/MamangRust/simple_microservice_ecommerce/user/internal/grpc_client"
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/handler"
	mencache "github.com/MamangRust/simple_microservice_ecommerce/user/internal/redis"
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/repository"
	"github.com/MamangRust/simple_microservice_ecommerce/user/internal/service"
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/database"
	db "github.com/MamangRust/simple_microservice_ecommerce/user/pkg/database/schema"
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/database/seeder"
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/dotenv"
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/hash"
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/logger"
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/otel"
	pbuser "github.com/MamangRust/simple_microservice_ecommerce_pb/user"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port int
)

func init() {
	port = viper.GetInt("GRPC_USER_PORT")

	if port == 0 {
		port = 50052
	}

	flag.IntVar(&port, "port", port, "gRPC server port")
}

func getEnvOrDefault(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)

	if !exists {
		return defaultValue
	}

	return value
}

func loadServiceAddresses() *ServiceAddresses {
	return &ServiceAddresses{
		Role: getEnvOrDefault("GRPC_ROLE_ADDR", "localhost:50052"),
	}
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

func createServiceConnections(addresses *ServiceAddresses, logger logger.LoggerInterface) (grpcclient.ServiceConnections, error) {
	var connections grpcclient.ServiceConnections

	conns := map[string]*string{
		"Role": &addresses.Role,
	}

	for name, addr := range conns {
		conn, err := createConnection(*addr, name, logger)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to connect to %s service, continuing with nil connection", name), zap.Error(err))
			conn = nil
		}

		switch name {
		case "Role":
			connections.Role = conn
		}
	}

	return connections, nil
}

func closeConnections(conns grpcclient.ServiceConnections, logger logger.LoggerInterface) {
	connsMap := map[string]*grpc.ClientConn{
		"Role": conns.Role,
	}

	for name, conn := range connsMap {
		if conn != nil {
			if err := conn.Close(); err != nil {
				logger.Error(fmt.Sprintf("Failed to close %s connection", name), zap.Error(err))
			}
		}
	}
}

type Server struct {
	DB          *db.Queries
	Services    service.Service
	Hash        hash.HashPassword
	Handlers    handler.Handler
	Ctx         context.Context
	Logger      logger.LoggerInterface
	Connections grpcclient.ServiceConnections
}

type ServiceAddresses struct {
	Role string
}

func NewServer(ctx context.Context) (*Server, func(context.Context) error, error) {
	flag.Parse()

	telemetry := otel.NewTelemetry(otel.Config{
		ServiceName:            "user-service",
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

	logger, err := logger.NewLogger("user-service", telemetry.GetLogger())
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

	hash := hash.NewHashingPassword()
	logger.Info("Password hasher initialized")

	repositories := repository.NewRepositories(DB)
	logger.Info("Repository layer initialized")

	myredis := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", viper.GetString("REDIS_HOST_USER"), viper.GetString("REDIS_PORT_USER")),
		Password:     viper.GetString("REDIS_PASSWORD_USER"),
		DB:           viper.GetInt("REDIS_DB_USER"),
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

	addresses := loadServiceAddresses()

	logger.Info("Loaded service addresses")

	connections, err := createServiceConnections(addresses, logger)
	if err != nil {
		logger.Error("Failed to create service connections, continuing with available connections", zap.Error(err))
	}
	logger.Info("Created gRPC service connections")

	client := grpcclient.NewHandlerGrpcClient(&grpcclient.Deps{
		ServiceConnections: &connections,
		Logger:             logger,
	})

	grpcClients := service.GrpcClient{
		RoleClient:     client.RoleClient,
		UserRoleClient: client.UserRoleClient,
	}

	services := service.NewService(&service.Deps{
		Reposotories: repositories,
		Hash:         hash,
		GrpcClient:   grpcClients,
		Logger:       logger,
		Mencache:     mencache,
	})

	logger.Info("Service layer initialized")

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  logger,
	})
	logger.Info("Handler layer initialized")

	logger.Info("Server dependencies initialized successfully")

	return &Server{
		DB:          DB,
		Services:    services,
		Handlers:    handlers,
		Ctx:         ctx,
		Hash:        hash,
		Logger:      logger,
		Connections: connections,
	}, telemetry.Shutdown, nil
}

func (s *Server) Run() {
	s.Logger.Info("Starting gRPC server...")

	defer func() {
		s.Logger.Info("Shutting down server, closing connections...")
		closeConnections(s.Connections, s.Logger)
		s.Logger.Info("All connections closed.")
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		s.Logger.Fatal("Failed to start TCP listener", zap.Error(err), zap.Int("port", port))
	}
	s.Logger.Info("TCP listener started successfully", zap.String("address", lis.Addr().String()))

	grpcServer := grpc.NewServer()

	s.Logger.Info("Registering gRPC services...")
	pbuser.RegisterUserQueryServiceServer(grpcServer, s.Handlers)
	pbuser.RegisterUserCommandServiceServer(grpcServer, s.Handlers)
	s.Logger.Info("gRPC services registered successfully")

	var wg sync.WaitGroup

	wg.Go(func() {
		s.Logger.Info("gRPC server is starting to serve requests", zap.Int("port", port))
		if err := grpcServer.Serve(lis); err != nil {
			s.Logger.Fatal("Failed to start gRPC server", zap.Error(err))
		}
	})

	wg.Wait()
}

func (s *Server) Seed() {
	dbSeederEnabled := viper.GetString("DB_SEEDER")
	s.Logger.Info("Checking database seeder flag", zap.String("DB_SEEDER", dbSeederEnabled))

	if dbSeederEnabled == "true" {
		s.Logger.Info("Database seeder is enabled, starting seeding process...")

		seeder := seeder.NewUserSeeder(s.DB, s.Hash, s.Logger)

		if err := seeder.SeedAll(s.Ctx); err != nil {
			s.Logger.Error("Failed to seed database", zap.Error(err))
			return
		}

		s.Logger.Info("Database seeding completed successfully.")
	} else {
		s.Logger.Info("Database seeder is disabled.")
	}
}
