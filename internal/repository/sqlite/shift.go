package sqlite

import (
	"context"
	"fmt"
	"time"

	"github.com/uchupx/kajian-api/pkg/db"

	"github.com/uchupx/worker-roster-management-system/internal/model"
	"github.com/uchupx/worker-roster-management-system/internal/repository"
)

type ShiftRepository struct {
	db *db.DB
}

func (s *ShiftRepository) CreateShift(ctx context.Context, shift model.Shift) (*int64, error) {
	stmt, err := s.db.FPreparexContext(ctx, CreateShiftQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	defer stmt.Close()

	res, err := stmt.ExecContext(ctx, shift.UserID, shift.StartTime, shift.EndTime, shift.ShiftDate, shift.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to execute statement: %w", err)
	}

	shiftID, err := res.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return &shiftID, nil
}

func (s *ShiftRepository) UpdateShift(ctx context.Context, shift model.Shift) error {
	stmt, err := s.db.FPreparexContext(ctx, UpdateShiftQuery)
	if err != nil {
		return fmt.Errorf("failed to prepare statement: %w", err)
	}

	defer stmt.Close()

	_, err = stmt.ExecContext(ctx, shift.StartTime, shift.EndTime, shift.ShiftDate, shift.Status, shift.ID)
	if err != nil {
		return fmt.Errorf("failed to execute statement: %w", err)
	}

	return nil
}

func (s ShiftRepository) FindShift(ctx context.Context, userID, status, id *int) ([]model.ShiftUser, error) {
	stmt, err := s.db.FPreparexContext(ctx, FindShiftQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	defer stmt.Close()

	var shifts []model.ShiftUser

	rows, err := stmt.QueryxContext(ctx, id, id, status, status, userID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query row: %w", err)
	}

	for rows.Next() {
		shift := model.ShiftUser{}
		err := rows.Scan(&shift.ID, &shift.UserID, &shift.StartTime, &shift.EndTime, &shift.ShiftDate, &shift.Status, &shift.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		shifts = append(shifts, shift)
	}

	return shifts, nil
}

func (s ShiftRepository) FindShiftByDate(ctx context.Context, startDate, endDate time.Time, approvedStatus *bool) ([]model.ShiftUser, error) {
	stmt, err := s.db.FPreparexContext(ctx, FindShiftByDateQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	defer stmt.Close()

	var shifts []model.ShiftUser
	rows, err := stmt.QueryxContext(ctx, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query row: %w", err)
	}

	for rows.Next() {
		shift := model.ShiftUser{}
		err := rows.Scan(&shift.ID, &shift.UserID, &shift.StartTime, &shift.EndTime, &shift.ShiftDate, &shift.Status, &shift.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		shifts = append(shifts, shift)
	}

	return shifts, nil
}

func (s ShiftRepository) FindShiftByUserIDDate(ctx context.Context, userID *int, startDate, endDate time.Time, approvedStatus *int8) ([]model.ShiftUser, error) {
	stmt, err := s.db.FPreparexContext(ctx, FindShiftByUserIDDateQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	defer stmt.Close()

	var shifts []model.ShiftUser
	rows, err := stmt.QueryxContext(ctx, userID, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query row: %w", err)
	}

	for rows.Next() {
		shift := model.ShiftUser{}
		err := rows.Scan(&shift.ID, &shift.UserID, &shift.StartTime, &shift.EndTime, &shift.ShiftDate, &shift.Status, &shift.Name)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		shifts = append(shifts, shift)
	}

	return shifts, nil
}

func (s ShiftRepository) FindShiftByID(ctx context.Context, userID int) (*model.Shift, error) {
	stmt, err := s.db.FPreparexContext(ctx, FindShiftByIDQuery)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare statement: %w", err)
	}

	defer stmt.Close()

	row := stmt.QueryRowxContext(ctx, userID)

	res := model.Shift{}
	if err = row.Scan(&res.ID, &res.UserID, &res.StartTime, &res.EndTime, &res.ShiftDate, &res.Status); err != nil {
		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	return &res, nil
}

func NewShiftRepository(db *db.DB) repository.ShiftRepositoryInterface {
	return &ShiftRepository{
		db: db,
	}
}
