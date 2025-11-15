package apps

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"

	grpcclient "github.com/MamangRust/simple_microservice_ecommerce/order/internal/grpc_client"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/handler"
	"go.uber.org/zap"

	pborder "github.com/MamangRust/simple_microservice_ecommerce_pb/order"
	pborderitem "github.com/MamangRust/simple_microservice_ecommerce_pb/order_item"

	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/repository"
	"github.com/MamangRust/simple_microservice_ecommerce/order/internal/service"
	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/database"
	db "github.com/MamangRust/simple_microservice_ecommerce/order/pkg/database/schema"
	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/database/seeder"
	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/dotenv"
	"github.com/MamangRust/simple_microservice_ecommerce/order/pkg/logger"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	port int
)

func init() {
	port = viper.GetInt("GRPC_ORDER_ADDR")
	if port == 0 {
		port = 50055
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
		Product: getEnvOrDefault("GRPC_PRODUCT_ADDR", "localhost:50054"),
		User:    getEnvOrDefault("GRPC_USER_ADDR", "localhost:50052"),
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
		"Product": &addresses.Product,
		"User":    &addresses.User,
	}

	for name, addr := range conns {
		conn, err := createConnection(*addr, name, logger)
		if err != nil {
			logger.Error(fmt.Sprintf("Failed to connect to %s service, continuing with nil connection", name), zap.Error(err))
			conn = nil
		}

		switch name {
		case "Product":
			connections.Product = conn
		case "User":
			connections.User = conn
		}
	}

	return connections, nil
}

func closeConnections(conns grpcclient.ServiceConnections, logger logger.LoggerInterface) {
	connsMap := map[string]*grpc.ClientConn{
		"Product": conns.Product,
		"User":    conns.User,
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
	Handlers    handler.Handler
	Ctx         context.Context
	Logger      logger.LoggerInterface
	Connections grpcclient.ServiceConnections
}

type ServiceAddresses struct {
	Product string
	User    string
}

func NewServer(ctx context.Context) (*Server, error) {
	flag.Parse()

	logger, err := logger.NewLogger("order-service")
	if err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	if err := dotenv.Viper(); err != nil {
		logger.Fatal("Failed to load .env file", zap.Error(err))
	}
	logger.Info("Successfully loaded .env file")

	conn, err := database.NewClient(logger)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	logger.Info("Successfully connected to the database")
	DB := db.New(conn)

	repositories := repository.NewRepository(DB)

	addresses := loadServiceAddresses()
	logger.Info("Loaded service addresses",
		zap.String("product_service_address", addresses.Product),
		zap.String("user_service_address", addresses.User),
	)

	connections, err := createServiceConnections(addresses, logger)
	if err != nil {
		logger.Error("Failed to create service connections, continuing with available connections", zap.Error(err))
	}
	logger.Info("Created gRPC service connections")

	client := grpcclient.NewHandlerGrpcClient(&connections, logger)

	grpcClients := service.GrpcClient{
		UserClient:    client.UserClient,
		ProductClient: client.ProductClient,
	}

	services := service.NewService(&service.Deps{
		Repository: repositories,
		GrpcClient: grpcClients,
		Logger:     logger,
	})
	logger.Info("Service layer initialized")

	handlers := handler.NewHandler(services, logger)
	logger.Info("Handler layer initialized")

	logger.Info("Server dependencies initialized successfully")
	return &Server{
		DB:          DB,
		Services:    services,
		Logger:      logger,
		Handlers:    handlers,
		Ctx:         ctx,
		Connections: connections,
	}, nil
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
	pborder.RegisterOrderQueryServiceServer(grpcServer, s.Handlers.OrderHandler())
	pborder.RegisterOrderCommandServiceServer(grpcServer, s.Handlers.OrderHandler())
	pborderitem.RegisterOrderItemServiceServer(grpcServer, s.Handlers.OrderItemHandler())
	s.Logger.Info("gRPC services registered successfully")

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		s.Logger.Info("gRPC server is starting to serve requests")
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

		seeder := seeder.NewOrderSeeder(s.DB, s.Logger)

		if err := seeder.SeedAll(s.Ctx); err != nil {
			s.Logger.Error("Failed to seed database", zap.Error(err))
			return
		}

		s.Logger.Info("Database seeding completed successfully.")
	} else {
		s.Logger.Info("Database seeder is disabled.")
	}
}
