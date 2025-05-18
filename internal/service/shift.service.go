package service

import (
	"context"
	"fmt"
	"time"

	api "github.com/uchupx/worker-roster-management-system/generated"
	"github.com/uchupx/worker-roster-management-system/internal/enums"
	"github.com/uchupx/worker-roster-management-system/internal/model"
	"github.com/uchupx/worker-roster-management-system/internal/repository"
	"github.com/uchupx/worker-roster-management-system/pkg/helper"
)

type ShiftService struct {
	userService UserService
	roleService RoleService

	shiftRepository repository.ShiftRepositoryInterface
}

type ShiftServiceParams struct {
	UserService     UserService
	ShiftRepository repository.ShiftRepositoryInterface
	RoleService     RoleService
}

func NewShiftService(params ShiftServiceParams) *ShiftService {
	return &ShiftService{
		userService:     params.UserService,
		shiftRepository: params.ShiftRepository,
		roleService:     params.RoleService,
	}
}

func GetFirstAndLastDayOfWeek(date time.Time) (time.Time, time.Time) {
	// Get the weekday (0 = Sunday, 1 = Monday, ..., 6 = Saturday)
	weekday := int(date.Weekday())

	// Calculate the first day of the week (assuming Sunday as the first day of the week)
	firstDay := date.AddDate(0, 0, -weekday).Truncate(24 * time.Hour)

	// Calculate the last day of the week
	lastDay := firstDay.AddDate(0, 0, 6).Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	return firstDay, lastDay
}

func (s *ShiftService) FindShift(ctx context.Context, params api.GetShiftsParams) ([]api.ShiftSchema, error) {
	shifts, err := s.shiftRepository.FindShift(ctx, params.UserId, params.Status, params.Id)
	if err != nil {
		return nil, helper.NewHTTPError(500, fmt.Errorf("failed to find shift by user ID: %w", err))
	}

	res := make([]api.ShiftSchema, len(shifts))

	for i, shift := range shifts {
		res[i] = *s.mapToShiftSchemaWithUser(&shift)
	}

	return res, nil
}

func (s *ShiftService) FindShiftByDate(ctx context.Context, dateStart, dateEnd time.Time) ([]api.ShiftSchema, error) {
	shifts, err := s.shiftRepository.FindShiftByDate(ctx, dateStart, dateEnd, nil)
	if err != nil {
		return nil, helper.NewHTTPError(500, fmt.Errorf("failed to find shift by date: %w", err))
	}

	res := make([]api.ShiftSchema, len(shifts))

	for i, shift := range shifts {
		res[i] = *s.mapToShiftSchemaWithUser(&shift)

	}

	return res, nil
}

func (s *ShiftService) CreateShift(ctx context.Context, req api.ShiftRequest) (*api.ShiftSchema, error) {
	if err := s.validateShiftExistence(ctx, req); err != nil {
		return nil, err
	}

	if err := s.validateWeeklyShiftLimit(ctx, req); err != nil {
		return nil, err
	}

	shift, err := s.createNewShift(ctx, req)
	if err != nil {
		return nil, err
	}

	return s.mapToShiftSchema(shift), nil
}

func (s *ShiftService) validateShiftExistence(ctx context.Context, req api.ShiftRequest) error {
	isExist, err := s.shiftRepository.FindShiftByUserIDDate(ctx, nil, req.ShiftDate, req.ShiftDate, nil)
	if err != nil {
		return helper.NewHTTPError(500, fmt.Errorf("failed to find shift by user ID and date: %w", err))
	}
	for _, shift := range isExist {
		fmt.Println(shift)
		if shift.Status == enums.ShiftStatusApproved {
			return helper.NewHTTPError(403, fmt.Errorf("shift already assigned"))
		}

		if req.UserId != nil && shift.UserID == int64(*req.UserId) && shift.Status == enums.ShiftStatusPending {
			return helper.NewHTTPError(403, fmt.Errorf("shift already requested"))
		}
	}

	return nil
}

func (s *ShiftService) validateWeeklyShiftLimit(ctx context.Context, req api.ShiftRequest) error {
	firstDayOfWeek, lastDayOfWeek := GetFirstAndLastDayOfWeek(req.ShiftDate)
	isExist, err := s.shiftRepository.FindShiftByUserIDDate(ctx, req.UserId, firstDayOfWeek, lastDayOfWeek, &enums.ShiftStatusApproved)
	if err != nil {
		return helper.NewHTTPError(500, fmt.Errorf("failed to find shift by user ID and date: %w", err))
	}

	if len(isExist) >= 5 {
		return helper.NewHTTPError(403, fmt.Errorf("user ID %d already has 5 shifts in the week of date %s", *req.UserId, req.ShiftDate))
	}

	return nil
}

func (s *ShiftService) createNewShift(ctx context.Context, req api.ShiftRequest) (*model.Shift, error) {
	shift := model.Shift{
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		ShiftDate: req.ShiftDate,
		UserID:    int64(*req.UserId),
	}

	authUserId := ctx.Value("userID")

	if authUserId != nil && authUserId.(int64) != shift.UserID {
		isAdmin, err := s.roleService.IsAdmin(ctx, int(authUserId.(int64)))
		if err != nil {
			return nil, err
		}

		if isAdmin {
			shift.Status = enums.ShiftStatusApproved
		}
	}

	shiftId, err := s.shiftRepository.CreateShift(ctx, shift)
	if err != nil {
		return nil, helper.NewHTTPError(500, fmt.Errorf("failed to create shift: %w", err))
	}

	shift.ID = *shiftId
	return &shift, nil
}

func (s *ShiftService) ApproveShift(ctx context.Context, shiftId, userId int) (*api.ShiftSchema, error) {
	isAdmin, err := s.roleService.IsAdmin(ctx, userId)

	if err != nil {
		return nil, helper.NewHTTPError(500, fmt.Errorf("failed to check admin role: %w", err))
	}

	if !isAdmin {
		return nil, helper.NewHTTPError(403, fmt.Errorf("user is not an admin"))
	}

	shift, err := s.getShiftById(ctx, shiftId)
	if err != nil {
		return nil, err
	}

	if err := s.validateShiftExistence(ctx, api.ShiftRequest{
		ShiftDate: shift.ShiftDate,
	}); err != nil {
		return nil, err
	}

	if shift.Status != enums.ShiftStatusPending {
		return nil, helper.NewHTTPError(403, fmt.Errorf("shift is not pending"))
	}

	shift.Status = enums.ShiftStatusApproved
	if err := s.updateShift(ctx, shift); err != nil {
		return nil, err
	}

	return s.mapToShiftSchema(shift), nil
}

func (s *ShiftService) RejectShift(ctx context.Context, shiftId, userId int) (*api.ShiftSchema, error) {
	isAdmin, err := s.roleService.IsAdmin(ctx, userId)
	if err != nil {
		return nil, helper.NewHTTPError(500, fmt.Errorf("failed to check admin role: %w", err))
	}

	if !isAdmin {
		return nil, helper.NewHTTPError(403, fmt.Errorf("user is not an admin"))
	}

	shift, err := s.getShiftById(ctx, shiftId)
	if err != nil {
		return nil, err
	}

	if shift.Status != enums.ShiftStatusPending {
		return nil, helper.NewHTTPError(403, fmt.Errorf("shift is not pending"))
	}

	shift.Status = enums.ShiftStatusRejected
	if err := s.updateShift(ctx, shift); err != nil {
		return nil, err
	}

	return s.mapToShiftSchema(shift), nil
}

func (s *ShiftService) getShiftById(ctx context.Context, shiftId int) (*model.Shift, error) {
	shift, err := s.shiftRepository.FindShiftByID(ctx, shiftId)
	if err != nil {
		return nil, helper.NewHTTPError(404, fmt.Errorf("shift not found: %w", err))
	}
	return shift, nil
}

func (s *ShiftService) updateShift(ctx context.Context, shift *model.Shift) error {
	if err := s.shiftRepository.UpdateShift(ctx, *shift); err != nil {
		return helper.NewHTTPError(500, fmt.Errorf("failed to update shift: %w", err))
	}
	return nil
}

func (s *ShiftService) mapToShiftSchemaWithUser(shift *model.ShiftUser) *api.ShiftSchema {
	return &api.ShiftSchema{
		Id:        int(shift.ID),
		UserId:    int(shift.UserID),
		UserName:  &shift.Name,
		StartTime: shift.StartTime,
		EndTime:   shift.EndTime,
		ShiftDate: shift.ShiftDate,
		Status:    intPointer(int(shift.Status)),
		CreatedAt: shift.CreatedAt,
	}
}

func (s *ShiftService) mapToShiftSchema(shift *model.Shift) *api.ShiftSchema {
	return &api.ShiftSchema{
		Id:        int(shift.ID),
		UserId:    int(shift.UserID),
		StartTime: shift.StartTime,
		EndTime:   shift.EndTime,
		ShiftDate: shift.ShiftDate,
		Status:    intPointer(int(shift.Status)),
		CreatedAt: shift.CreatedAt,
	}
}

func intPointer(val int) *int {
	return &val
}
