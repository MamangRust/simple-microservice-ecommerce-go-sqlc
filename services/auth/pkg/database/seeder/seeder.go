package seeder

import (
	"context"
	"fmt"
	"time"

	db "github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/database/schema"
	"github.com/MamangRust/simple_microservice_ecommerce/auth/pkg/logger"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthSeeder struct {
	db     *db.Queries
	logger logger.LoggerInterface
}

func NewAuthSeeder(db *db.Queries, logger logger.LoggerInterface) *AuthSeeder {
	return &AuthSeeder{
		db:     db,
		logger: logger,
	}
}

func (s *AuthSeeder) SeedAll(ctx context.Context) error {
	s.logger.Info("[SEEDER] Starting auth service seeding...")

	if err := s.SeedRefreshTokens(ctx); err != nil {
		s.logger.Error("[SEEDER] Failed to seed refresh tokens", zap.Error(err))
		return fmt.Errorf("failed to seed refresh tokens: %w", err)
	}

	if err := s.SeedResetTokens(ctx); err != nil {
		s.logger.Error("[SEEDER] Failed to seed reset tokens", zap.Error(err))
		return fmt.Errorf("failed to seed reset tokens: %w", err)
	}

	s.logger.Info("[SEEDER] Auth service seeding completed successfully")
	return nil
}

func (s *AuthSeeder) SeedRefreshTokens(ctx context.Context) error {
	s.logger.Info("[SEEDER] Seeding refresh tokens...")

	seedData := []db.CreateRefreshTokenParams{
		{
			UserID:     1,
			Token:      uuid.New().String(),
			Expiration: time.Now().Add(1 * time.Hour),
		},
		{
			UserID:     2,
			Token:      uuid.New().String(),
			Expiration: time.Now().Add(1 * time.Hour),
		},
	}

	createdCount := 0
	for _, data := range seedData {
		existingToken, err := s.db.FindRefreshTokenByUserId(ctx, int32(data.UserID))
		if err == nil && existingToken != nil {
			s.logger.Info("[SEEDER] Refresh token already exists, skipping", zap.Int64("user_id", int64(data.UserID)))
			continue
		}

		_, err = s.db.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
			UserID:     int32(data.UserID),
			Token:      data.Token,
			Expiration: data.Expiration,
		})

		if err != nil {
			s.logger.Error("[SEEDER] Error inserting refresh token", zap.Error(err), zap.Int64("user_id", int64(data.UserID)))
			return err
		}

		s.logger.Info("[SEEDER] Successfully created refresh token", zap.Int64("user_id", int64(data.UserID)))
		createdCount++
	}

	s.logger.Info("[SEEDER] Refresh token seeding summary", zap.Int("total_tokens_in_list", len(seedData)), zap.Int("tokens_created", createdCount))
	return nil
}

func (s *AuthSeeder) SeedResetTokens(ctx context.Context) error {
	s.logger.Info("[SEEDER] Seeding reset tokens...")

	seedData := []db.CreateResetTokenParams{
		{
			UserID:     1,
			Token:      uuid.New().String(),
			ExpiryDate: time.Now().Add(1 * time.Hour),
		},
		{
			UserID:     2,
			Token:      uuid.New().String(),
			ExpiryDate: time.Now().Add(1 * time.Hour),
		},
	}

	createdCount := 0
	for _, data := range seedData {
		existingToken, err := s.db.GetResetToken(ctx, data.Token)
		if err == nil && existingToken != nil {
			s.logger.Info("[SEEDER] Reset token already exists, skipping", zap.Int64("user_id", data.UserID))
			continue
		}

		_, err = s.db.CreateResetToken(ctx, db.CreateResetTokenParams{
			UserID:     data.UserID,
			Token:      data.Token,
			ExpiryDate: data.ExpiryDate,
		})

		if err != nil {
			s.logger.Error("[SEEDER] Error inserting reset token", zap.Error(err), zap.Int64("user_id", data.UserID))
			return err
		}

		s.logger.Info("[SEEDER] Successfully created reset token", zap.Int64("user_id", data.UserID))
		createdCount++
	}

	s.logger.Info("[SEEDER] Reset token seeding summary", zap.Int("total_tokens_in_list", len(seedData)), zap.Int("tokens_created", createdCount))
	return nil
}

func (s *AuthSeeder) ClearAll(ctx context.Context) error {
	s.logger.Info("[SEEDER] Clearing auth service data...")

	deletedRefreshTokenCount := 0
	for _, userID := range []int32{1, 2, 3} {
		err := s.db.DeleteRefreshTokenByUserId(ctx, userID)
		if err != nil {
			s.logger.Error("[SEEDER] Error clearing refresh tokens", zap.Error(err), zap.Int32("user_id", userID))
			return err
		}
		deletedRefreshTokenCount++
	}

	deletedResetTokenCount := 0
	for _, userID := range []int64{1, 2} {
		err := s.db.DeleteResetToken(ctx, userID)
		if err != nil {
			s.logger.Error("[SEEDER] Error clearing reset token", zap.Error(err), zap.Int64("user_id", userID))
			return err
		}
		deletedResetTokenCount++
	}

	s.logger.Info("[SEEDER] Auth service data cleared successfully",
		zap.Int("deleted_refresh_tokens", deletedRefreshTokenCount),
		zap.Int("deleted_reset_tokens", deletedResetTokenCount),
	)
	return nil
}
