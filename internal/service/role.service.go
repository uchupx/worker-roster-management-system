package service

import (
	"context"
	"fmt"

	"github.com/uchupx/worker-roster-management-system/internal/model"
	"github.com/uchupx/worker-roster-management-system/internal/repository"
	"github.com/uchupx/worker-roster-management-system/pkg/helper"
)


type RoleService struct {
	roleRepository repository.RoleRepositoryInterface
}

type RoleServiceParams struct {
	RoleRepository repository.RoleRepositoryInterface
}

func NewRoleService(params RoleServiceParams) *RoleService {
	return &RoleService{
		roleRepository: params.RoleRepository,
	}
}


func (s *RoleService) GetAllUserRole(ctx context.Context, userId int) ([]model.Role, error) {
    // Fetch roles for the given user ID
    roles, err := s.roleRepository.FindByUserId(ctx, userId)
    if err != nil {
        return nil, helper.NewHTTPError(500, fmt.Errorf("failed to find roles by user ID: %w", err))
    }
    return roles, nil
}

func (s *RoleService) GetRoleById(ctx context.Context, id int) (*model.Role, error) {
    // Fetch a role by its ID
    role, err := s.roleRepository.FindById(ctx, id)
    if err != nil {
        return nil, helper.NewHTTPError(404, fmt.Errorf("role not found for ID %d: %w", id, err))
    }
    return role, nil
}

func (s *RoleService) IsAdmin(ctx context.Context, userId int) (bool, error) {
    // Fetch roles for the given user ID
    roles, err := s.roleRepository.FindByUserId(ctx, userId)
    if err != nil {
        return false, helper.NewHTTPError(500, fmt.Errorf("failed to find roles by user ID: %w", err))
    }

    // Check if the user has an admin role
    return s.hasAdminRole(roles), nil
}

// Helper function to check if a user has an admin role
func (s *RoleService) hasAdminRole(roles []model.Role) bool {
    for _, role := range roles {
        if role.ID == int64(1) { // Assuming role ID 1 represents admin
            return true
        }
    }
    return false
}
