package repository

import (
	refreshtokenrepository "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/repository/refresh_token"
	resettokenrepository "github.com/MamangRust/simple_microservice_ecommerce/auth/internal/repository/reset_token"
	db "github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/database/schema"
)

type Repository interface {
	refreshtokenrepository.RefreshTokenRepository
	resettokenrepository.ResetTokenRepository
}

type repository struct {
	refreshtokenrepository.RefreshTokenRepository
	resettokenrepository.ResetTokenRepository
}

func NewRepository(
	db *db.Queries,
) Repository {
	return &repository{
		RefreshTokenRepository: refreshtokenrepository.NewRefreshTokenRepository(db),
		ResetTokenRepository:   resettokenrepository.NewResetTokenRepository(db),
	}
}
