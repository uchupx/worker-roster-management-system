package repository

import (
	"context"
	"time"

	"github.com/uchupx/worker-roster-management-system/internal/model"
)

type ShiftRepositoryInterface interface {
	CreateShift(ctx context.Context, shift model.Shift) (*int64, error)
	UpdateShift(ctx context.Context, shift model.Shift) error

	FindShiftByID(ctx context.Context, userID int) (*model.Shift, error)
	FindShift(ctx context.Context, userID, status, id *int) ([]model.ShiftUser, error)
	FindShiftByUserIDDate(ctx context.Context, userID *int, startDate, endDate time.Time, approvedStatus *int8) ([]model.ShiftUser, error)
	FindShiftByDate(ctx context.Context, startDate, endDate time.Time, approvedStatus *bool) ([]model.ShiftUser, error)
}
