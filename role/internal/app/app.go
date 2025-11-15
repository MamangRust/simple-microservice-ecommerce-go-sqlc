package apps

import (
	"context"
	"flag"
	"fmt"
	"net"
	"sync"

	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/handler"
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/repository"
	"github.com/MamangRust/simple_microservice_ecommerce/role/internal/service"
	"github.com/MamangRust/simple_microservice_ecommerce/role/pkg/database"
	db "github.com/MamangRust/simple_microservice_ecommerce/role/pkg/database/schema"
	"github.com/MamangRust/simple_microservice_ecommerce/role/pkg/database/seeder"
	"github.com/MamangRust/simple_microservice_ecommerce/role/pkg/dotenv"
	"github.com/MamangRust/simple_microservice_ecommerce/role/pkg/logger"
	pbrole "github.com/MamangRust/simple_microservice_ecommerce_pb/role"
	pbuserrole "github.com/MamangRust/simple_microservice_ecommerce_pb/user_role"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

var (
	port int
)

func init() {
	port = viper.GetInt("GRPC_ROLE_PORT")

	if port == 0 {
		port = 50053
	}

	flag.IntVar(&port, "port", port, "gRPC server port")
}

type Server struct {
	DB       *db.Queries
	Logger   logger.LoggerInterface
	Services service.Service
	Handlers handler.Handler
	Ctx      context.Context
}

func NewServer(ctx context.Context) (*Server, error) {
	flag.Parse()

	logger, err := logger.NewLogger("role-service")
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

	repositories := repository.NewRepositories(DB)
	logger.Info("Repository layer initialized")

	services := service.NewService(repositories, logger)
	logger.Info("Service layer initialized")

	handlers := handler.NewHandler(services, logger)
	logger.Info("Handler layer initialized")

	logger.Info("Server dependencies initialized successfully")
	return &Server{
		DB:       DB,
		Services: services,
		Handlers: handlers,
		Ctx:      ctx,
		Logger:   logger,
	}, nil
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
	pbrole.RegisterRoleQueryServiceServer(grpcServer, s.Handlers.RoleHandler())
	pbrole.RegisterRoleCommandServiceServer(grpcServer, s.Handlers.RoleHandler())
	pbuserrole.RegisterUserRoleServiceServer(grpcServer, s.Handlers.UserRoleHandler())
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

		seeder := seeder.NewRoleSeeder(s.DB, s.Logger)

		if err := seeder.SeedAll(s.Ctx); err != nil {
			s.Logger.Error("Failed to seed database", zap.Error(err))
			return
		}

		s.Logger.Info("Database seeding completed successfully.")
	} else {
		s.Logger.Info("Database seeder is disabled.")
	}
}
