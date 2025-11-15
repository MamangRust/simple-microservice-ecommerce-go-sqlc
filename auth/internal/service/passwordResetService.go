package service

import (
	"context"
	"encoding/json"
	"strconv"
	"time"

	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/requests"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/response"
	grpcclient "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/grpc_client"
	resettokenrepository "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/repository/reset_token"
	emails "github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/email"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/kafka"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/logger"
	randomstring "github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/random_string"
	"go.uber.org/zap"
)

type passwordResetServiceDeps struct {
	kafka      *kafka.Kafka
	UserClient grpcclient.UserGrpcClientHandler
	ResetToken resettokenrepository.ResetTokenRepository
	logger     logger.LoggerInterface
}

type passwordResetService struct {
	kafka      *kafka.Kafka
	userClient grpcclient.UserGrpcClientHandler
	resetToken resettokenrepository.ResetTokenRepository
	logger     logger.LoggerInterface
}

func NewPasswordResetService(params *passwordResetServiceDeps) PasswordResetService {
	return &passwordResetService{
		kafka:      params.kafka,
		userClient: params.UserClient,
		resetToken: params.ResetToken,
		logger:     params.logger,
	}
}

func (s *passwordResetService) ForgotPassword(ctx context.Context, email string) (bool, *response.ErrorResponse) {
	s.logger.Info("Forgot password request", zap.String("email", email))

	res, errResp := s.userClient.FindByEmail(ctx, email)

	if errResp != nil {
		s.logger.Error("User not found for forgot password", zap.String("email", email), zap.Any("error", errResp))
		return false, response.NewErrorResponse("user not found", 404)
	}

	random, err := randomstring.GenerateRandomString(10)

	if err != nil {
		s.logger.Error("Failed to generate reset token", zap.String("email", email), zap.Error(err))
		return false, response.NewErrorResponse("failed generate random string", 400)
	}

	_, err = s.resetToken.CreateResetToken(ctx, &requests.CreateResetTokenRequest{
		UserID:     res.Data.ID,
		ResetToken: random,
		ExpiredAt:  time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
	})

	if err != nil {
		s.logger.Error("Failed to create reset token", zap.Int("user_id", res.Data.ID), zap.Error(err))
		return false, response.NewErrorResponse("failed create reset token", 400)
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
		s.logger.Error("Failed to marshal email payload", zap.String("email", email), zap.Error(err))
		return false, response.NewErrorResponse("failed to marshal email payload", 400)
	}

	err = s.kafka.SendMessage("email-service-topic-auth-forgot-password", strconv.Itoa(res.Data.ID), payloadBytes)
	if err != nil {
		s.logger.Error("Failed to send forgot password email", zap.String("email", email), zap.Int("user_id", res.Data.ID), zap.Error(err))
		return false, response.NewErrorResponse("failed to send email message", 400)
	}

	s.logger.Info("Forgot password email sent successfully", zap.String("email", email), zap.Int("user_id", res.Data.ID))

	return true, nil
}

func (s *passwordResetService) ResetPassword(ctx context.Context, request *requests.CreateResetPasswordRequest) (bool, *response.ErrorResponse) {
	s.logger.Info("Password reset request", zap.String("reset_token", request.ResetToken))

	res, err := s.resetToken.FindByResetToken(ctx, request.ResetToken)

	var userid int

	if err != nil {
		s.logger.Error("Invalid or expired reset token", zap.String("reset_token", request.ResetToken), zap.Error(err))
		return false, response.NewErrorResponse("invalid or expired reset token", 400)
	}

	userid = int(res.UserID)

	if request.Password != request.ConfirmPassword {
		s.logger.Error("Password confirmation mismatch", zap.Int("user_id", userid))
		return false, response.NewErrorResponse("password no match", 400)
	}

	_, errResp := s.userClient.UpdateUserPassword(ctx, &requests.UpdateUserPasswordRequest{
		UserID:   userid,
		Password: request.Password,
	})

	if errResp != nil {
		s.logger.Error("Failed to update user password", zap.Int("user_id", userid), zap.Any("error", errResp))
		return false, response.NewErrorResponse("failed to update", 400)
	}

	s.logger.Info("Password reset successful", zap.Int("user_id", userid))

	return true, nil
}

func (s *passwordResetService) VerifyCode(ctx context.Context, code string) (bool, *response.ErrorResponse) {
	s.logger.Info("Verification code request", zap.String("code", code))

	res, errResp := s.userClient.FindVerificationCode(ctx, code)

	if errResp != nil {
		s.logger.Error("Invalid verification code", zap.String("code", code), zap.Any("error", errResp))
		return false, response.NewErrorResponse("invalid verification code", 400)
	}

	_, errResp = s.userClient.UpdateUserIsVerified(ctx, &requests.UpdateUserVerifiedRequest{
		UserID:     res.Data.ID,
		IsVerified: true,
	})

	if errResp != nil {
		s.logger.Error("Failed to update user verification", zap.Int("user_id", res.Data.ID), zap.Any("error", errResp))
		return false, response.NewErrorResponse("failed update user verify", 400)
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
		s.logger.Error("Failed to marshal verification email payload", zap.Int("user_id", res.Data.ID), zap.Error(err))
		return false, response.NewErrorResponse("failed to marshal verification email payload", 400)
	}

	err = s.kafka.SendMessage("email-service-topic-auth-verify-code-success", strconv.Itoa(res.Data.ID), payloadBytes)
	if err != nil {
		s.logger.Error("Failed to send verification success email", zap.Int("user_id", res.Data.ID), zap.Error(err))
		return false, response.NewErrorResponse("failed to send verification email", 400)
	}

	s.logger.Info("User verification successful", zap.Int("user_id", res.Data.ID), zap.String("email", res.Data.Email))

	return true, nil
}
