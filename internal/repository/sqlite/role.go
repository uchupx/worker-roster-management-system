package sqlite

import (
	"context"
	"fmt"

	"github.com/uchupx/kajian-api/pkg/db"
	"github.com/uchupx/worker-roster-management-system/internal/model"
	"github.com/uchupx/worker-roster-management-system/internal/repository"
)


var (
	findRoleByIdQuery = `SELECT roles.id, roles.name, roles.created_at FROM roles where roles.id = ?`

	findRoleByUserIdQuery = `SELECT roles.id, roles.name, roles.created_at FROM roles JOIN user_roles on user_roles.role_id = roles.id where user_roles.user_id = ?`
)


type RoleRepository struct {
	db *db.DB
}


func (r *RoleRepository) FindById(ctx context.Context, id int) (role *model.Role,err error) {
	stmt, err := r.db.FPreparexContext(ctx, findRoleByIdQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	defer stmt.Close()

	row := stmt.QueryRowContext(ctx, id)

	err = row.Scan(&role.ID, &role.Name, &role.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	return role, nil
}


func (r *RoleRepository) FindByUserId(ctx context.Context, userId int) (roles []model.Role, err error) {
	stmt, err := r.db.FPreparexContext(ctx, findRoleByUserIdQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to query rows: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var role model.Role
		err = rows.Scan(&role.ID, &role.Name, &role.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func NewRoleRepository(db *db.DB) repository.RoleRepositoryInterface {
	return &RoleRepository{
		db: db,
	}
}
