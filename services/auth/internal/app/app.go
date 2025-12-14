package apps

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	grpcclient "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/grpc_client"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/handler"
	mencache "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/redis"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/repository"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/service"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/auth"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/database"
	db "github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/database/schema"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/database/seeder"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/dotenv"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/hash"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/kafka"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/logger"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/otel"
	pbauth "github.com/MamangRust/simple_microservice_ecommerce_pb/auth"
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
	port = viper.GetInt("GRPC_AUTH_ADDR")
	if port == 0 {
		port = 50051
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
		User: getEnvOrDefault("GRPC_USER_ADDR", "localhost:50053"),
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
		"User": &addresses.User,
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
		case "User":
			connections.User = conn
		}
	}

	return connections, nil
}

func closeConnections(conns grpcclient.ServiceConnections, logger logger.LoggerInterface) {
	connsMap := map[string]*grpc.ClientConn{
		"Role": conns.Role,
		"User": conns.User,
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
	DB           *db.Queries
	TokenManager *auth.Manager
	Logger       logger.LoggerInterface
	Services     service.Service
	Handlers     handler.Handler
	Ctx          context.Context
	Connections  grpcclient.ServiceConnections
}

type ServiceAddresses struct {
	Role string
	User string
}

func NewServer(ctx context.Context) (*Server, func(context.Context) error, error) {
	flag.Parse()

	telemetry := otel.NewTelemetry(otel.Config{
		ServiceName:            "auth-service",
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

	logger, err := logger.NewLogger("auth-service", telemetry.GetLogger())
	if err != nil {
		return nil, nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	if err := dotenv.Viper(); err != nil {
		logger.Fatal("Failed to load .env file", zap.Error(err))
	}

	tokenManager, err := auth.NewManager(viper.GetString("SECRET_KEY"))
	if err != nil {
		logger.Fatal("Failed to create token manager", zap.Error(err))
	}

	conn, err := database.NewClient(ctx)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	DB := db.New(conn)

	hash := hash.NewHashingPassword()
	repositories := repository.NewRepository(DB)

	myredis := redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%s", viper.GetString("REDIS_HOST_AUTH"), viper.GetString("REDIS_PORT")),
		Password:     viper.GetString("REDIS_PASSWORD_AUTH"),
		DB:           viper.GetInt("REDIS_DB_AUTH"),
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

	kafka := kafka.NewKafka([]string{viper.GetString("KAFKA_BROKERS")})

	addresses := loadServiceAddresses()
	connections, err := createServiceConnections(addresses, logger)
	if err != nil {
		logger.Error("Failed to create service connections, continuing with available connections", zap.Error(err))
	}

	client := grpcclient.NewHandlerGrpcClient(&grpcclient.Deps{
		ServiceConnections: &connections,
		Logger:             logger,
	})

	grpcClients := service.GrpcClient{
		UserClient:     client.UserClient,
		RoleClient:     client.RoleClient,
		UserRoleClient: client.UserRoleClient,
	}

	services := service.NewService(&service.Deps{
		Repository: repositories,
		Hash:       hash,
		Token:      tokenManager,
		Kafka:      kafka,
		GrpcClient: grpcClients,
		Logger:     logger,
		Mencache:   mencache,
	})

	handlers := handler.NewHandler(&handler.Deps{
		Service: services,
		Logger:  logger,
	})

	return &Server{
		DB:           DB,
		TokenManager: tokenManager,
		Services:     services,
		Handlers:     handlers,
		Ctx:          ctx,
		Logger:       logger,
		Connections:  connections,
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
		s.Logger.Fatal("Failed to listen", zap.Error(err))
	}

	grpcServer := grpc.NewServer()
	pbauth.RegisterAuthServiceServer(grpcServer, s.Handlers)

	s.Logger.Info("gRPC server listening", zap.Int("port", port))

	if err := grpcServer.Serve(lis); err != nil {
		s.Logger.Fatal("Failed to start gRPC server", zap.Error(err))
	}
}

func (s *Server) Seed() {
	dbSeeder := viper.GetString("DB_SEEDER")

	if dbSeeder == "true" {
		s.Logger.Info("[SEEDER] Database seeder enabled, starting seeding process...")

		seeder := seeder.NewAuthSeeder(s.DB, s.Logger)

		if err := seeder.SeedAll(s.Ctx); err != nil {
			s.Logger.Error("[SEEDER] Failed to seed database", zap.Error(err))
			return
		}

		s.Logger.Info("[SEEDER] Database seeding completed successfully.")
	}
}
