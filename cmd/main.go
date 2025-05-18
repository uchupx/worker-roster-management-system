package main

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/uchupx/worker-roster-management-system/config"
	api "github.com/uchupx/worker-roster-management-system/generated"
	"github.com/uchupx/worker-roster-management-system/handler"
	"github.com/uchupx/worker-roster-management-system/internal/database"
	myMiddleware "github.com/uchupx/worker-roster-management-system/internal/middleware"
	"github.com/uchupx/worker-roster-management-system/pkg/helper"
)

func main() {

	conf := config.GetConfig()

	e := echo.New()
	var server api.ServerInterface = newServer(conf)
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	registerHandlers(e, conf, server)

	e.Use(middleware.Logger())

	e.Logger.Fatal(e.Start(":" + conf.App.Port))
}

func registerHandlers(e *echo.Echo, conf *config.Config, server api.ServerInterface) {
	swagger, err := api.GetSwagger()
	if err != nil {
		panic(err)
	}

	middle := myMiddleware.New(myMiddleware.Config{
		Redis: database.GetRedisClient(*conf),
	})

	e.POST("/login", server.PostLogin)
	e.GET("/docs", func(c echo.Context) error {
		return c.JSON(http.StatusOK, swagger)
	})

	auth := e.Group("", middle.Authorization)
	auth.POST("/register", server.PostRegister)

	auth.GET("/me/shifts", server.GetMeShifts)
	auth.POST("/me/shifts", server.PostMeShifts)

	auth.POST("/shifts", server.PostShifts)
	auth.GET("/shifts", func(c echo.Context) error {
		var status *int
		var id *int

		if c.QueryParam("status") != "" {
			statusTemp := int(helper.StringToInt(c.QueryParam("status")))
			status = &statusTemp
		}

		if c.QueryParam("id") != "" {
			shiftIdTemp := int(helper.StringToInt(c.QueryParam("id")))
			id = &shiftIdTemp
		}

		return server.GetShifts(c, api.GetShiftsParams{
			Status: status,
			Id:     id,
		})
	})

	auth.GET("/shifts/months", func(c echo.Context) error {
		start := c.QueryParam("dateStart")
		end := c.QueryParam("dateEnd")

		startTime := helper.StringToInt(start)
		endTime := helper.StringToInt(end)

		return server.GetShiftsMonths(c, api.GetShiftsMonthsParams{
			DateStart: int(startTime),
			DateEnd:   int(endTime),
		})
	})

	auth.POST("/shifts/:id/approve", func(c echo.Context) error {
		id := c.Param("id")

		return server.PostShiftsIdApprove(c, id)
	})

	auth.POST("/shifts/:id/reject", func(c echo.Context) error {
		id := c.Param("id")

		return server.PostShiftsIdReject(c, id)
	})

	auth.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "pong")
	})

	auth.GET("/users", func(c echo.Context) error {
		var role *int

		roleStr := c.QueryParam("role")
		if roleStr != "" {
			roleTemp := int(helper.StringToInt(roleStr))
			role = &roleTemp
		}

		return server.GetUsers(c, api.GetUsersParams{
			Role: role,
		})
	})

	auth.POST("/users", server.PostUsers)

	auth.GET("/me", server.GetMe)
}

func newServer(conf *config.Config) *handler.Server {
	return handler.NewServer(conf)
}
