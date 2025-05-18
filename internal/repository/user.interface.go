package repository

import (
	"context"

	"github.com/uchupx/worker-roster-management-system/internal/model"
)

type UserRepositoryInterface interface {
	Insert(ctx context.Context, user *model.User, roleId int) (*int64, error)
	FindByEmail(ctx context.Context, username string) (*model.User, error)
	FindById(ctx context.Context, id int) (*model.User, error)
	FindUser(ctx context.Context, role *int, email *string, id *int) ([]model.User, error)
	UpdateUserRole(ctx context.Context, userId, roleId int) error
}
