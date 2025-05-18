package handler

import (
	echo "github.com/labstack/echo/v4"

	api "github.com/uchupx/worker-roster-management-system/generated"
	"github.com/uchupx/worker-roster-management-system/pkg/helper"
)

func (Server) responseError(ctx echo.Context, err error) error {
	if httpErr, ok := err.(*helper.HTTPError); ok {
		return ctx.JSON(httpErr.Code, api.ErrorResponse{
			Message: err.Error(),
		})
	}

	return ctx.JSON(500, api.ErrorResponse{
		Message: err.Error(),
	})
}
