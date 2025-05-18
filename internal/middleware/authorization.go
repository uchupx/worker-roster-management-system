package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	echo "github.com/labstack/echo/v4"
	redis "github.com/redis/go-redis/v9"

	"github.com/uchupx/worker-roster-management-system/internal/enums"
	"github.com/uchupx/worker-roster-management-system/internal/model"
)

func (m *Middleware) Authorization(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		auth := c.Request().Header.Get("Authorization")
		reg := regexp.MustCompile(`Bearer[\s]`)

		auth = reg.ReplaceAllString(auth, "")
		if strings.TrimSpace(auth) != "" {
			// Get the value from Redis
			result, err := m.Redis.Get(c.Request().Context(), fmt.Sprintf(enums.RedisKeyAuthorization, auth)).Result()
			if err == redis.Nil {
				return echo.NewHTTPError(401, "Unauthorized")
			} else if err != nil {
				return echo.NewHTTPError(500, "Internal Server Error")
			}

			var user model.User
			if err := json.Unmarshal([]byte(result), &user); err != nil {
				return echo.NewHTTPError(500, "Failed to parse user data")
			}
			c.Set("userID", user.ID)
			// Store authData in the context
			ctx := context.WithValue(c.Request().Context(), "authData", user)
			c.SetRequest(c.Request().WithContext(ctx))

			// Proceed to the next handler
			return next(c)
		} else {
			return echo.NewHTTPError(401, "Unauthorized")
		}
	}
}
