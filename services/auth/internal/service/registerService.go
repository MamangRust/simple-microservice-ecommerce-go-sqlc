package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/response"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/errorhandler"
	grpcclient "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/grpc_client"
	mencache "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/redis"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/email"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/kafka"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/logger"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/observability"
	randomstring "github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/random_string"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type RegisterServiceDeps struct {
	userClient        grpcclient.UserGrpcClientHandler
	roleClient        grpcclient.RoleGrpcClientHandler
	userRoleClient    grpcclient.UserRoleGrpcClientHandler
	kafka             *kafka.Kafka
	errorRandomString errorhandler.RandomStringErrorHandler
	errorMarshal      errorhandler.MarshalErrorHandler
	errorKafka        errorhandler.KafkaErrorHandler
	errorhandler      errorhandler.RegisterErrorHandler
	mencache          mencache.RegisterCache
	logger            logger.LoggerInterface
}

type registerService struct {
	userClient        grpcclient.UserGrpcClientHandler
	roleClient        grpcclient.RoleGrpcClientHandler
	userRoleClient    grpcclient.UserRoleGrpcClientHandler
	kafka             *kafka.Kafka
	logger            logger.LoggerInterface
	errorRandomString errorhandler.RandomStringErrorHandler
	errorMarshal      errorhandler.MarshalErrorHandler
	errorKafka        errorhandler.KafkaErrorHandler
	errorhandler      errorhandler.RegisterErrorHandler
	mencache          mencache.RegisterCache
	observability     observability.TraceLoggerObservability
}

func NewRegisterService(params *RegisterServiceDeps) *registerService {
	observability, _ := observability.NewObservability("register-service", params.logger)

	return &registerService{
		userClient:        params.userClient,
		roleClient:        params.roleClient,
		userRoleClient:    params.userRoleClient,
		kafka:             params.kafka,
		logger:            params.logger,
		errorRandomString: params.errorRandomString,
		errorMarshal:      params.errorMarshal,
		errorKafka:        params.errorKafka,
		errorhandler:      params.errorhandler,
		observability:     observability,
		mencache:          params.mencache,
	}
}

func (s *registerService) Register(ctx context.Context, request *requests.RegisterRequest) (*response.UserResponse, *response.ErrorResponse) {
	const method = "Register"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("email", request.Email))

	defer func() {
		end(status)
	}()

	random, err := randomstring.GenerateRandomString(10)
	if err != nil {
		status = "error"

		return s.errorRandomString.HandleRandomStringErrorRegister(err, "Register", "REGISTER_ERR", span, &status, zap.Error(err))
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
		status = "error"

		return s.errorhandler.HandleCreateUserError(err, "Register", "REGISTER_ERR", span, &status, zap.Error(err))
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
		status = "error"

		return s.errorMarshal.HandleMarshalRegisterError(err, "Register", "MARSHAL_ERR", span, &status, zap.Error(err))
	}

	err = s.kafka.SendMessage("email-service-topic-auth-register", strconv.Itoa(newUser.Data.ID), payloadBytes)

	if err != nil {
		status = "error"

		return s.errorKafka.HandleSendEmailRegister(err, "Register", "SEND_EMAIL_ERR", span, &status, zap.Error(err))
	}

	logSuccess("User registered successfully",
		zap.String("email", request.Email),
		zap.String("first_name", request.FirstName),
		zap.String("last_name", request.LastName),
	)

	s.mencache.SetVerificationCodeCache(ctx, request.Email, random, 15*time.Minute)

	return newUser.Data, nil
}
