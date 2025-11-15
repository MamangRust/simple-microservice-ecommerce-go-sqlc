package seeder

import (
	"context"
	"fmt"

	db "github.com/MamangRust/simple_microservice_ecommerce/role/pkg/database/schema"
	"github.com/MamangRust/simple_microservice_ecommerce/role/pkg/logger"
	"go.uber.org/zap"
)

type RoleSeeder struct {
	db     *db.Queries
	logger logger.LoggerInterface
}

func NewRoleSeeder(db *db.Queries, logger logger.LoggerInterface) *RoleSeeder {
	return &RoleSeeder{
		db:     db,
		logger: logger,
	}
}

func (s *RoleSeeder) SeedAll(ctx context.Context) error {
	s.logger.Info("[SEEDER] Starting role service seeding...")

	if err := s.SeedRoles(ctx); err != nil {
		s.logger.Error("[SEEDER] Failed to seed roles", zap.Error(err))
		return fmt.Errorf("failed to seed roles: %w", err)
	}

	if err := s.SeedUserRoles(ctx); err != nil {
		s.logger.Error("[SEEDER] Failed to seed user roles", zap.Error(err))
		return fmt.Errorf("failed to seed user roles: %w", err)
	}

	s.logger.Info("[SEEDER] Role service seeding completed successfully")
	return nil
}

func (s *RoleSeeder) SeedRoles(ctx context.Context) error {
	s.logger.Info("[SEEDER] Seeding roles...")

	roles := []string{
		"admin",
		"user",
		"moderator",
		"vendor",
		"customer",
	}

	createdCount := 0
	for _, roleName := range roles {
		existingRole, err := s.db.GetRoleByName(ctx, roleName)
		if err == nil && existingRole != nil {
			s.logger.Info("[SEEDER] Role already exists, skipping", zap.String("role_name", roleName))
			continue
		}

		role, err := s.db.CreateRole(ctx, roleName)
		if err != nil {
			s.logger.Error("[SEEDER] Error inserting role", zap.String("role_name", roleName), zap.Error(err))
			return err
		}

		s.logger.Info("[SEEDER] Successfully created role", zap.String("role_name", roleName), zap.Int32("role_id", role.RoleID))
		createdCount++
	}

	s.logger.Info("[SEEDER] Role seeding summary", zap.Int("total_roles_in_list", len(roles)), zap.Int("roles_created", createdCount))
	return nil
}

func (s *RoleSeeder) SeedUserRoles(ctx context.Context) error {
	s.logger.Info("[SEEDER] Seeding user roles...")

	userRoles := []struct {
		UserID int32
		RoleID int32
	}{
		{1, 1},
		{2, 2},
		{3, 2},
		{4, 3},
		{5, 4},
		{1, 2},
	}

	createdCount := 0
	for _, ur := range userRoles {
		userExistingRoles, err := s.db.GetUserRoles(ctx, ur.UserID)
		if err == nil {
			roleExists := false
			for _, existingRole := range userExistingRoles {
				if existingRole.RoleID == ur.RoleID {
					roleExists = true
					break
				}
			}
			if roleExists {
				s.logger.Info("[SEEDER] User role assignment already exists, skipping",
					zap.Int32("user_id", ur.UserID),
					zap.Int32("role_id", ur.RoleID),
				)
				continue
			}
		}

		_, err = s.db.AssignRoleToUser(ctx, db.AssignRoleToUserParams{
			UserID: ur.UserID,
			RoleID: ur.RoleID,
		})

		if err != nil {
			s.logger.Error("[SEEDER] Error assigning role to user",
				zap.Error(err),
				zap.Int32("user_id", ur.UserID),
				zap.Int32("role_id", ur.RoleID),
			)
			return err
		}

		s.logger.Info("[SEEDER] Successfully assigned role to user",
			zap.Int32("user_id", ur.UserID),
			zap.Int32("role_id", ur.RoleID),
		)
		createdCount++
	}

	s.logger.Info("[SEEDER] User role seeding summary", zap.Int("total_assignments_in_list", len(userRoles)), zap.Int("assignments_created", createdCount))
	return nil
}

func (s *RoleSeeder) ClearAll(ctx context.Context) error {
	s.logger.Info("[SEEDER] Clearing role service data...")

	deletedUserRoleAssignments := 0

	for userID := int32(1); userID <= 5; userID++ {
		roles, err := s.db.GetUserRoles(ctx, userID)
		if err == nil {
			for _, role := range roles {
				err = s.db.RemoveRoleFromUser(ctx, db.RemoveRoleFromUserParams{
					UserID: userID,
					RoleID: role.RoleID,
				})
				if err != nil {
					s.logger.Error("[SEEDER] Error removing role from user",
						zap.Error(err),
						zap.Int32("user_id", userID),
						zap.Int32("role_id", role.RoleID),
					)
					return err
				}
				deletedUserRoleAssignments++
			}
		}
	}

	deletedRoles := 0
	roles := []string{"admin", "user", "moderator", "vendor", "customer"}
	for _, roleName := range roles {
		role, err := s.db.GetRoleByName(ctx, roleName)
		if err == nil && role != nil {
			err = s.db.DeletePermanentRole(ctx, role.RoleID)
			if err != nil {
				s.logger.Error("[SEEDER] Error deleting role", zap.String("role_name", roleName), zap.Error(err))
				return err
			}
			deletedRoles++
		}
	}

	s.logger.Info("[SEEDER] Role service data cleared successfully",
		zap.Int("deleted_user_role_assignments", deletedUserRoleAssignments),
		zap.Int("deleted_roles", deletedRoles),
	)

	return nil
}
