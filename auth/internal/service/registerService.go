package service

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/response"
	grpcclient "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/grpc_client"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/email"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/kafka"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/logger"
	randomstring "github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/random_string"
	"go.uber.org/zap"
)

type RegisterServiceDeps struct {
	userClient     grpcclient.UserGrpcClientHandler
	roleClient     grpcclient.RoleGrpcClientHandler
	userRoleClient grpcclient.UserRoleGrpcClientHandler
	kafka          *kafka.Kafka
	logger         logger.LoggerInterface
}

type registerService struct {
	userClient     grpcclient.UserGrpcClientHandler
	roleClient     grpcclient.RoleGrpcClientHandler
	userRoleClient grpcclient.UserRoleGrpcClientHandler
	kafka          *kafka.Kafka
	logger         logger.LoggerInterface
}

func NewRegisterService(params *RegisterServiceDeps) *registerService {
	return &registerService{
		userClient:     params.userClient,
		roleClient:     params.roleClient,
		userRoleClient: params.userRoleClient,
		kafka:          params.kafka,
		logger:         params.logger,
	}
}

func (s *registerService) Register(ctx context.Context, request *requests.RegisterRequest) (*response.UserResponse, *response.ErrorResponse) {
	s.logger.Info("Registration attempt", zap.String("email", request.Email), zap.String("first_name", request.FirstName), zap.String("last_name", request.LastName))

	random, err := randomstring.GenerateRandomString(10)
	if err != nil {
		s.logger.Error("Failed to generate verification code", zap.String("email", request.Email), zap.Error(err))
		return nil, response.NewErrorResponse("failed to generate verification code", 400)
	}

	newUser, errResp := s.userClient.CreateUser(ctx, &requests.CreateUserRequest{
		FirstName:       request.FirstName,
		LastName:        request.LastName,
		Email:           request.Email,
		Password:        request.Password,
		ConfirmPassword: request.ConfirmPassword,
		VerifiedCode:    random,
		IsVerified:      false,
	})

	if errResp != nil {
		s.logger.Error("Failed to create user", zap.String("email", request.Email), zap.Any("error", errResp))
		return nil, response.NewErrorResponse("failed to create user", 400)
	}

	htmlBody := email.GenerateEmailHTML(map[string]string{
		"Title":   "Welcome to SanEdge",
		"Message": "Your account has been successfully created.",
		"Button":  "Verify Now",
		"Link":    "https://sanedge.example.com/login?verify_code=" + random,
	})

	emailPayload := map[string]any{
		"email":   request.Email,
		"subject": "Welcome to SanEdge",
		"body":    htmlBody,
	}

	payloadBytes, err := json.Marshal(emailPayload)

	if err != nil {
		s.logger.Error("Failed to encode email payload", zap.String("email", request.Email), zap.Error(err))
		return nil, response.NewErrorResponse("failed to encode email payload", 400)
	}

	err = s.kafka.SendMessage("email-service-topic-auth-register", strconv.Itoa(newUser.Data.ID), payloadBytes)

	if err != nil {
		s.logger.Error("Failed to send email event", zap.String("email", request.Email), zap.Int("user_id", newUser.Data.ID), zap.Error(err))
		return nil, response.NewErrorResponse("failed to send email event", 400)
	}

	s.logger.Info("User registration successful", zap.Int("user_id", newUser.Data.ID), zap.String("email", request.Email))

	return newUser.Data, nil
}
