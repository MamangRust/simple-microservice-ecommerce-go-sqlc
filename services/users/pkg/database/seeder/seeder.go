package seeder

import (
	"context"
	"fmt"

	db "github.com/MamangRust/simple_microservice_ecommerce/user/pkg/database/schema"
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/hash"
	"github.com/MamangRust/simple_microservice_ecommerce/user/pkg/logger"
	"go.uber.org/zap"
)

type UserSeeder struct {
	db     *db.Queries
	hash   hash.HashPassword
	logger logger.LoggerInterface
}

func NewUserSeeder(db *db.Queries, hash hash.HashPassword, logger logger.LoggerInterface) *UserSeeder {
	return &UserSeeder{
		db:     db,
		hash:   hash,
		logger: logger,
	}
}

func (s *UserSeeder) SeedAll(ctx context.Context) error {
	s.logger.Info("[SEEDER] Starting users service seeding...")

	if err := s.SeedUsers(ctx); err != nil {
		s.logger.Error("[SEEDER] Failed to seed users", zap.Error(err))
		return fmt.Errorf("failed to seed users: %w", err)
	}

	s.logger.Info("[SEEDER] Users service seeding completed successfully")
	return nil
}

func (s *UserSeeder) SeedUsers(ctx context.Context) error {
	s.logger.Info("[SEEDER] Seeding users...")

	users := []db.CreateUserParams{
		{
			Firstname: "Admin",
			Lastname:  "User",
			Email:     "admin@example.com",
			Password:  "admin123",
		},
		{
			Firstname: "John",
			Lastname:  "Doe",
			Email:     "john.doe@example.com",
			Password:  "password123",
		},
		{
			Firstname: "Jane",
			Lastname:  "Smith",
			Email:     "jane.smith@example.com",
			Password:  "password123",
		},
		{
			Firstname: "Moderator",
			Lastname:  "User",
			Email:     "moderator@example.com",
			Password:  "mod123",
		},
		{
			Firstname: "Vendor",
			Lastname:  "Account",
			Email:     "vendor@example.com",
			Password:  "vendor123",
		},
		{
			Firstname: "Unverified",
			Lastname:  "User",
			Email:     "unverified@example.com",
			Password:  "temp123",
		},
	}

	createdCount := 0
	for _, user := range users {
		existingUser, err := s.db.GetUserByEmail(ctx, user.Email)
		if err == nil && existingUser != nil {
			s.logger.Info("[SEEDER] User already exists, skipping", zap.String("email", user.Email))
			continue
		}

		hashedPassword, err := s.hash.HashPassword(user.Password)
		if err != nil {
			s.logger.Error("[SEEDER] Error hashing password for user", zap.String("email", user.Email), zap.Error(err))
			return err
		}

		verified := true

		newUser, err := s.db.CreateUser(ctx, db.CreateUserParams{
			Firstname:        user.Firstname,
			Lastname:         user.Lastname,
			Email:            user.Email,
			Password:         string(hashedPassword),
			VerificationCode: "hello",
			IsVerified:       &verified,
		})

		if err != nil {
			s.logger.Error("[SEEDER] Error inserting user", zap.String("email", user.Email), zap.Error(err))
			return err
		}

		s.logger.Info("[SEEDER] Successfully created user", zap.String("email", user.Email), zap.Int("user_id", int(newUser.UserID)))
		createdCount++
	}

	s.logger.Info("[SEEDER] User seeding summary", zap.Int("total_users_in_list", len(users)), zap.Int("users_created", createdCount))
	return nil
}

func (s *UserSeeder) ClearAll(ctx context.Context) error {
	s.logger.Info("[SEEDER] Clearing users service data...")

	users, err := s.db.GetUsers(ctx, db.GetUsersParams{
		Column1: "",
		Limit:   1000,
		Offset:  0,
	})

	if err == nil {
		deletedCount := 0
		for _, user := range users {
			err = s.db.DeleteUserPermanently(ctx, user.UserID)
			if err != nil {
				s.logger.Error("[SEEDER] Error deleting user", zap.Int("user_id", int(user.UserID)), zap.Error(err))
				return err
			}
			deletedCount++
		}
		s.logger.Info("[SEEDER] User clearing summary", zap.Int("users_found", len(users)), zap.Int("users_deleted", deletedCount))
	} else {
		s.logger.Error("[SEEDER] Failed to get users list for clearing", zap.Error(err))
	}

	s.logger.Info("[SEEDER] Users service data cleared successfully")
	return nil
}
