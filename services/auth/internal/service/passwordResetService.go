package service

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/response"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/errorhandler"
	grpcclient "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/grpc_client"
	mencache "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/redis"
	resettokenrepository "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/repository/reset_token"
	emails "github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/email"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/kafka"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/logger"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/observability"
	randomstring "github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/random_string"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

type passwordResetServiceDeps struct {
	kafka             *kafka.Kafka
	UserClient        grpcclient.UserGrpcClientHandler
	ResetToken        resettokenrepository.ResetTokenRepository
	logger            logger.LoggerInterface
	errorhandler      errorhandler.PasswordResetErrorHandler
	errorRandomString errorhandler.RandomStringErrorHandler
	errorMarshal      errorhandler.MarshalErrorHandler
	errorPassword     errorhandler.PasswordErrorHandler
	errorKafka        errorhandler.KafkaErrorHandler
	mencache          mencache.PasswordResetCache
}

type passwordResetService struct {
	kafka             *kafka.Kafka
	userClient        grpcclient.UserGrpcClientHandler
	resetToken        resettokenrepository.ResetTokenRepository
	logger            logger.LoggerInterface
	errorhandler      errorhandler.PasswordResetErrorHandler
	errorRandomString errorhandler.RandomStringErrorHandler
	errorMarshal      errorhandler.MarshalErrorHandler
	errorPassword     errorhandler.PasswordErrorHandler
	errorKafka        errorhandler.KafkaErrorHandler
	observability     observability.TraceLoggerObservability
	mencache          mencache.PasswordResetCache
}

func NewPasswordResetService(params *passwordResetServiceDeps) PasswordResetService {
	observability, _ := observability.NewObservability("password-service", params.logger)

	return &passwordResetService{
		kafka:             params.kafka,
		userClient:        params.UserClient,
		resetToken:        params.ResetToken,
		logger:            params.logger,
		errorhandler:      params.errorhandler,
		errorRandomString: params.errorRandomString,
		errorMarshal:      params.errorMarshal,
		errorPassword:     params.errorPassword,
		errorKafka:        params.errorKafka,
		observability:     observability,
		mencache:          params.mencache,
	}
}

func (s *passwordResetService) ForgotPassword(ctx context.Context, email string) (bool, *response.ErrorResponse) {
	const method = "ForgotPassword"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("email", email))

	defer func() {
		end(status)
	}()

	res, errResp := s.userClient.FindByEmail(ctx, email)

	if errResp != nil {
		status = "error"

		return s.errorhandler.HandleFindEmailError(errResp.ToGRPCError(), method, "FORGOT_PASSWORD_ERR", span, &status, zap.String("email", email))
	}

	random, err := randomstring.GenerateRandomString(10)

	if err != nil {
		status = "error"

		return s.errorRandomString.HandleRandomStringErrorForgotPassword(err, method, "FORGOT_PASSWORD_ERR", span, &status, zap.String("email", email), zap.Error(err))
	}

	_, err = s.resetToken.CreateResetToken(ctx, &requests.CreateResetTokenRequest{
		UserID:     res.Data.ID,
		ResetToken: random,
		ExpiredAt:  time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
	})

	if err != nil {
		status = "error"

		return s.errorhandler.HandleCreateResetTokenError(err, method, "FORGOT_PASSWORD_ERR", span, &status, zap.String("email", email), zap.Error(err))
	}

	htmlBody := emails.GenerateEmailHTML(map[string]string{
		"Title":   "Reset Your Password",
		"Message": "Click the button below to reset your password.",
		"Button":  "Reset Password",
		"Link":    "https://sanedge.example.com/reset-password?token=" + random,
	})

	emailPayload := map[string]any{
		"email":   res.Data.Email,
		"subject": "Password Reset Request",
		"body":    htmlBody,
	}

	payloadBytes, err := json.Marshal(emailPayload)
	if err != nil {
		status = "error"

		return s.errorMarshal.HandleMarsalForgotPassword(err, method, "FORGOT_PASSWORD_ERR", span, &status, zap.Error(err))
	}

	err = s.kafka.SendMessage("email-service-topic-auth-forgot-password", strconv.Itoa(res.Data.ID), payloadBytes)
	if err != nil {
		status = "error"

		return s.errorKafka.HandleSendEmailForgotPassword(err, method, "FORGOT_PASSWORD_ERR", span, &status, zap.Error(err))
	}

	logSuccess("Successfully sent password reset email", zap.String("email", email))

	return true, nil
}

func (s *passwordResetService) ResetPassword(ctx context.Context, request *requests.CreateResetPasswordRequest) (bool, *response.ErrorResponse) {
	const method = "ResetPassword"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("reset_token", request.ResetToken))

	defer func() {
		end(status)
	}()

	var userID int
	var found bool

	userID, found = s.mencache.GetResetTokenCache(ctx, request.ResetToken)
	if !found {
		res, err := s.resetToken.FindByResetToken(ctx, request.ResetToken)
		if err != nil {
			return s.errorhandler.HandleFindTokenError(err, method, "RESET_PASSWORD_ERR", span, &status, zap.String("reset_token", request.ResetToken))
		}
		userID = int(res.UserID)

		s.mencache.SetResetTokenCache(ctx, request.ResetToken, userID, 5*time.Minute)
	}

	if request.Password != request.ConfirmPassword {
		status = "error"

		err := errors.New("password and confirm password do not match")

		return s.errorPassword.HandlePasswordNotMatchError(err, method, "RESET_PASSWORD_ERR", span, &status, zap.String("reset_token", request.ResetToken))
	}

	_, errResp := s.userClient.UpdateUserPassword(ctx, &requests.UpdateUserPasswordRequest{
		UserID:   userID,
		Password: request.Password,
	})

	if errResp != nil {
		status = "error"

		return s.errorhandler.HandleUpdatePasswordError(errResp.ToGRPCError(), method, "RESET_PASSWORD_ERR", span, &status, zap.String("reset_token", request.ResetToken))
	}

	_ = s.resetToken.DeleteResetToken(ctx, userID)
	s.mencache.DeleteResetTokenCache(ctx, request.ResetToken)

	logSuccess("Successfully reset password", zap.String("reset_token", request.ResetToken))

	return true, nil
}

func (s *passwordResetService) VerifyCode(ctx context.Context, code string) (bool, *response.ErrorResponse) {
	const method = "VerifyCode"

	ctx, span, end, status, logSuccess := s.observability.StartTracingAndLogging(ctx, method, attribute.String("code", code))

	defer func() {
		end(status)
	}()

	res, errResp := s.userClient.FindVerificationCode(ctx, code)

	if errResp != nil {
		status = "error"

		return s.errorhandler.HandleVerifyCodeError(errResp.ToGRPCError(), method, "VERIFY_CODE_ERR", span, &status, zap.String("code", code))
	}

	_, errResp = s.userClient.UpdateUserIsVerified(ctx, &requests.UpdateUserVerifiedRequest{
		UserID:     res.Data.ID,
		IsVerified: true,
	})

	s.mencache.DeleteVerificationCodeCache(ctx, res.Data.Email)

	if errResp != nil {
		status = "error"

		return s.errorhandler.HandleUpdateVerifiedError(errResp.ToGRPCError(), method, "VERIFY_CODE_ERR", span, &status, zap.Int("user.id", res.Data.ID))
	}

	htmlBody := emails.GenerateEmailHTML(map[string]string{
		"Title":   "Verification Success",
		"Message": "Your account has been successfully verified. Click the button below to view or manage your card.",
		"Button":  "Go to Dashboard",
		"Link":    "https://sanedge.example.com/card/create",
	})

	emailPayload := map[string]any{
		"email":   res.Data.Email,
		"subject": "Verification Success",
		"body":    htmlBody,
	}

	payloadBytes, err := json.Marshal(emailPayload)
	if err != nil {
		status = "error"

		return s.errorMarshal.HandleMarshalVerifyCode(err, method, "SEND_EMAIL_VERIFY_CODE_ERR", span, &status, zap.Error(err))
	}

	err = s.kafka.SendMessage("email-service-topic-auth-verify-code-success", strconv.Itoa(res.Data.ID), payloadBytes)

	if err != nil {
		status = "error"

		return s.errorKafka.HandleSendEmailVerifyCode(err, method, "SEND_EMAIL_VERIFY_CODE_ERR", span, &status, zap.Error(err))
	}

	logSuccess("Successfully verify code", zap.String("code", code))

	return true, nil
}
