package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/MamangRust/simple_microservice_ecommerce/auth/internal/domain/requests"
	refreshtokenrepository "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/repository/refresh_token"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/auth"
)

type tokenServiceDeps struct {
	Token               auth.TokenManager
	RefreshTokenCommand refreshtokenrepository.RefreshTokenCommandRepository
}

type tokenService struct {
	refreshTokenCommand refreshtokenrepository.RefreshTokenCommandRepository
	token               auth.TokenManager
}

func NewTokenService(
	params *tokenServiceDeps,
) *tokenService {
	return &tokenService{
		refreshTokenCommand: params.RefreshTokenCommand,
		token:               params.Token,
	}
}

func (s *tokenService) createAccessToken(id int) (string, error) {
	res, err := s.token.GenerateToken(id, "access")

	if err != nil {
		return "", err
	}

	return res, nil
}

func (s *tokenService) createRefreshToken(ctx context.Context, id int) (string, error) {
	res, err := s.token.GenerateToken(id, "refresh")

	if err != nil {
		return "", err
	}

	if err := s.refreshTokenCommand.DeleteRefreshTokenByUserId(ctx, id); err != nil && !errors.Is(err, sql.ErrNoRows) {
		return "", err
	}

	_, err = s.refreshTokenCommand.CreateRefreshToken(ctx, &requests.CreateRefreshToken{
		Token:     res,
		UserId:    id,
		ExpiresAt: time.Now().Add(24 * time.Hour).Format("2006-01-02 15:04:05"),
	})
	if err != nil {
		return "", err
	}

	return res, nil
}
