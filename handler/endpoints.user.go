package handler

import (
	echo "github.com/labstack/echo/v4"

	api "github.com/uchupx/worker-roster-management-system/generated"
	"github.com/uchupx/worker-roster-management-system/pkg/helper"
)

func (s *Server) PostLogin(ctx echo.Context) error {
	var err error
	var body api.PostLoginJSONRequestBody

	if err = ctx.Bind(&body); err != nil {
		return s.responseError(ctx, helper.NewHTTPError(400, err))
	}

	res, err := s.authService.Login(ctx.Request().Context(), body)
	if err != nil {
		return s.responseError(ctx, err)
	}

	return ctx.JSON(200, res)
}

func (s *Server) PostRegister(ctx echo.Context) error {
	var err error
	var body api.PostRegisterJSONRequestBody

	if err = ctx.Bind(&body); err != nil {
		return s.responseError(ctx, helper.NewHTTPError(400, err))
	}

	res, err := s.userService.CreateUser(ctx.Request().Context(), body)
	if err != nil {
		return s.responseError(ctx, err)
	}

	return ctx.JSON(201, api.CreatedResponse{
		Id: string(res.ID),
	})
}


func (s *Server) GetMe(ctx echo.Context) error {
	userID := int(ctx.Get("userID").(int64))

	res, err := s.userService.GetUserById(ctx.Request().Context(), userID)
	if err != nil {
		return s.responseError(ctx, err)
	}

	return ctx.JSON(200, res)
}

func (s *Server) GetUsers(ctx echo.Context, params api.GetUsersParams) error {
	res, err := s.userService.GetUsers(ctx.Request().Context(), params.Role)
	if err != nil {
		return s.responseError(ctx, err)
	}

	return ctx.JSON(200, res)
}


func (s *Server) PostUsers(ctx echo.Context) error {
	var err error
	var body api.PostUsersJSONRequestBody

	if err = ctx.Bind(&body); err != nil {
		return s.responseError(ctx, helper.NewHTTPError(400, err))
	}

	res, err := s.userService.CreateUser(ctx.Request().Context(), body)
	if err != nil {
		return s.responseError(ctx, err)
	}

	return ctx.JSON(201, api.CreatedResponse{
		Id: string(res.ID),
	})
}
