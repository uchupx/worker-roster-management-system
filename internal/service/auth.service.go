package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	redis "github.com/redis/go-redis/v9"

	api "github.com/uchupx/worker-roster-management-system/generated"
	"github.com/uchupx/worker-roster-management-system/internal/enums"
	"github.com/uchupx/worker-roster-management-system/internal/repository"
	"github.com/uchupx/worker-roster-management-system/internal/service/jwt"
)

type AuthService struct {
	user        repository.UserRepositoryInterface
	crpyt       jwt.CryptService
	redisClient *redis.Client
}

type AuthServiceParams struct {
	User        repository.UserRepositoryInterface
	JWT         jwt.CryptService
	RedisClient *redis.Client
}

func NewAuthService(params AuthServiceParams) *AuthService {
	return &AuthService{
		user:        params.User,
		crpyt:       params.JWT,
		redisClient: params.RedisClient,
	}
}

func (s *AuthService) Login(ctx context.Context, req api.LoginRequest) (*api.LoginResponse, error) {
	user, err := s.user.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("error when find user by username: %w", err)
	} else if user == nil {
		return nil, enums.ErrNotFound
	}

	isValid, err := s.crpyt.Verify(req.Password, user.Password)
	if err != nil {
		return nil, fmt.Errorf("error when verify value: %w", err)
	}

	if !isValid {
		return nil, enums.ErrNotFound
	}

	tokenDuration := 1 * time.Hour
	token, err := s.crpyt.CreateAccessToken(tokenDuration, user)
	if err != nil {
		return nil, fmt.Errorf("error when create access token: %w", err)
	}

	if err := s.redisClient.Set(ctx, fmt.Sprintf(enums.RedisKeyAuthorization, *token), jsonStringify(user), tokenDuration).Err(); err != nil {
		return nil, fmt.Errorf("error when set redis: %w", err)
	}

	return &api.LoginResponse{
		Token:    *token,
		Duration: int(tokenDuration.Seconds()),
	}, nil
}

func jsonStringify(v interface{}) string {
	b, _ := json.Marshal(v)
	return string(b)
}
