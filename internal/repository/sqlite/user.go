package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/uchupx/kajian-api/pkg/db"

	"github.com/uchupx/worker-roster-management-system/internal/model"
	"github.com/uchupx/worker-roster-management-system/internal/repository"
)

type UserRepository struct {
	db *db.DB
}

func (r *UserRepository) Insert(ctx context.Context, data *model.User, roleId int) (*int64, error) {
	var lastInsertedId int64

	err := r.db.FTransaction(ctx, func(ctx context.Context, tx *sqlx.Tx) error {
		stmt, err := r.db.FPreparexContext(ctx, insertUserQuery)
		if err != nil {
			return fmt.Errorf("failed to prepare statement: %w", err)
		}

		defer stmt.Close()

		res, err := stmt.FExecContext(ctx,
			data.Name,
			data.Email,
			data.Password,
		)

		if err != nil {
			return fmt.Errorf("failed to execute statement: %w", err)
		}

		id, err := res.LastInsertId()

		if err != nil {
			return fmt.Errorf("failed to get last insert id: %w", err)
		}

		lastInsertedId = id

		if err = r.insertUserRole(ctx, int(lastInsertedId), roleId); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert user role: %w", err)
		}

		return nil
	}, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return &lastInsertedId, nil

}

func (r *UserRepository) FindUser(ctx context.Context, role *int, email *string, id *int) ([]model.User, error) {
	stmt, err := r.db.FPreparexContext(ctx, findUserQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	defer stmt.Close()

	var users []model.User
	rows, err := stmt.FQueryxContext(ctx, email, email, role, role, id, id)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}

	for rows.Next() {
		var user model.User

		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User

	stmt, err := r.db.FPreparexContext(ctx, findUserByEmailQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	defer stmt.Close()

	row := stmt.FQueryRowxContext(ctx, email)
	err = row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) FindById(ctx context.Context, id int) (*model.User, error) {
	stmt, err := r.db.FPreparexContext(ctx, findUserByIdQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	defer stmt.Close()

	data := &model.User{}

	row := stmt.FQueryRowxContext(ctx, id)

	err = row.Scan(
		&data.ID,
		&data.Name,
		&data.Email,
		&data.Password,
		&data.CreatedAt,
		&data.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	return data, nil
}

func (s *UserRepository) UpdateUserRole(ctx context.Context, userId, roleId int) error {
	// Implementation for updating a user's role in the database
	return nil
}

func (s UserRepository) insertUserRole(ctx context.Context, userId, roleId int) error {
	// Implementation for inserting a user role into the database
	return nil
}

func NewUserRepository(db *db.DB) repository.UserRepositoryInterface {
	return &UserRepository{
		db: db,
	}
}
