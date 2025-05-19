package handler

import (
	"context"
	"time"

	echo "github.com/labstack/echo/v4"

	api "github.com/uchupx/worker-roster-management-system/generated"
	"github.com/uchupx/worker-roster-management-system/pkg/helper"
)

func (s Server) GetMeShifts(ctx echo.Context) error {
	userID := int(ctx.Get("userID").(int64))

	shifts, err := s.shiftService.FindShift(ctx.Request().Context(), api.GetShiftsParams{UserId: &userID})
	if err != nil {
		return s.responseError(ctx, err)
	}

	return ctx.JSON(200, shifts)
}

func (s Server) PostMeShifts(ctx echo.Context) error {
	userID := int(ctx.Get("userID").(int64))
	var body api.ShiftRequest

	if err := ctx.Bind(&body); err != nil {
		return s.responseError(ctx, helper.NewHTTPError(400, err))
	}

	body.UserId = &userID

	shifts, err := s.shiftService.CreateShift(ctx.Request().Context(), body)
	if err != nil {
		return s.responseError(ctx, err)
	}

	return ctx.JSON(200, shifts)
}

func (s Server) GetShiftsMonths(ctx echo.Context, params api.GetShiftsMonthsParams) error {
	dateStart := time.Unix(int64(params.DateStart), 0).UTC()
	dateEnd := time.Unix(int64(params.DateEnd), 0).UTC()

	if params.DateStart == 0 || params.DateEnd == 0 {
		date := time.Now()
		dateStart = time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
		dateEnd = dateStart.AddDate(0, 1, 0).Add(-time.Second)
	}

	shifts, err := s.shiftService.FindShiftByDate(ctx.Request().Context(), dateStart, dateEnd)
	if err != nil {
		return s.responseError(ctx, err)
	}

	return ctx.JSON(200, shifts)
}

func (s Server) GetShifts(ctx echo.Context, params api.GetShiftsParams) error {
	shifts, err := s.shiftService.FindShift(ctx.Request().Context(), params)
	if err != nil {
		return s.responseError(ctx, err)
	}

	return ctx.JSON(200, shifts)
}

func (s Server) PostShifts(ctx echo.Context) error {
	userID := ctx.Get("userID").(int64)
	var body api.ShiftRequest

	if err := ctx.Bind(&body); err != nil {
		return s.responseError(ctx, helper.NewHTTPError(400, err))
	}

	c := ctx.Request().Context()
	c = context.WithValue(c, "userID", userID)

	shifts, err := s.shiftService.CreateShift(c, body)
	if err != nil {
		return s.responseError(ctx, err)
	}

	return ctx.JSON(200, shifts)
}

func (s Server) PostShiftsIdApprove(ctx echo.Context, id string) error {
	userID := int(ctx.Get("userID").(int64))

	_, err := s.shiftService.ApproveShift(ctx.Request().Context(), int(helper.StringToInt(id)), userID)
	if err != nil {
		return s.responseError(ctx, err)
	}

	return ctx.JSON(200, api.OkResponse{
		Id: id,
	})
}

func (s Server) PostShiftsIdReject(ctx echo.Context, id string) error {
	userID := int(ctx.Get("userID").(int64))

	_, err := s.shiftService.RejectShift(ctx.Request().Context(), int(helper.StringToInt(id)), userID)
	if err != nil {
		return s.responseError(ctx, err)
	}

	return ctx.JSON(200, api.OkResponse{
		Id: id,
	})
}
