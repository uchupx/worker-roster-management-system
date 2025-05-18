package repository

import (
	"context"

	"github.com/uchupx/worker-roster-management-system/internal/model"
)

type RoleRepositoryInterface interface {
	FindById(ctx context.Context, id int) (*model.Role, error)
	FindByUserId(ctx context.Context, userId int) ([]model.Role, error)
}
