package service

import (
	"context"
	"fmt"

	api "github.com/uchupx/worker-roster-management-system/generated"
	"github.com/uchupx/worker-roster-management-system/internal/model"
	"github.com/uchupx/worker-roster-management-system/internal/repository"
	"github.com/uchupx/worker-roster-management-system/internal/service/jwt"
	"github.com/uchupx/worker-roster-management-system/pkg/helper"
)

type UserService struct {
	user        repository.UserRepositoryInterface
	roleService RoleService
	crypt       jwt.CryptService
}

type UserServiceParams struct {
	User        repository.UserRepositoryInterface
	RoleService RoleService
	JWT         jwt.CryptService
}

func NewUserService(params UserServiceParams) *UserService {
	return &UserService{
		user:        params.User,
		crypt:       params.JWT,
		roleService: params.RoleService,
	}
}

func (s *UserService) GetUserById(ctx context.Context, id int) (*api.UserSchema, error) {
	fmt.Println("GetUserById", id)
	user, err := s.user.FindById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}

	if user == nil {
		return nil, helper.NewHTTPError(404, fmt.Errorf("user not found"))
	}

	roles, err := s.roleService.GetAllUserRole(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find user roles: %w", err)
	}

	if len(roles) == 0 {
		return nil, helper.NewHTTPError(404, fmt.Errorf("user roles not found"))
	}

	
	userWithRole := s.mapUserToUserRole(user, roles)
	return s.mapUserToUserSchema(userWithRole), nil
}

func (s *UserService) CreateUser(ctx context.Context, req api.RegisterRequest) (*model.User, error) {
	user := s.mapRequestToUser(req)

	if err := s.setSignedPassword(req.Password, user); err != nil {
		return nil, helper.NewHTTPError(400, fmt.Errorf("failed to set signed password: %w", err))
	}

	if err := s.insertUser(ctx, user, req.RoleId); err != nil {
		return nil, helper.NewHTTPError(500, fmt.Errorf("failed to insert user: %w", err))
	}

	return user, nil
}

func (s *UserService) mapRequestToUser(req api.RegisterRequest) *model.User {
	return &model.User{
		Name:  req.Name,
		Email: req.Email,
	}
}

func (s *UserService) setSignedPassword(password string, user *model.User) error {
	signPassword, err := s.crypt.CreateSignPSS(password)
	if err != nil {
		return fmt.Errorf("error creating signed password: %w", err)
	}
	user.Password = signPassword
	return nil
}

func (s *UserService) insertUser(ctx context.Context, user *model.User, roleId int) error {
	id, err := s.user.Insert(ctx, user, roleId)
	if err != nil {
		return fmt.Errorf("error inserting user: %w", err)
	}
	user.ID = *id
	return nil
}

func (s *UserService) GetUsers(ctx context.Context, roleId *int) ([]*api.UserSchema, error) {
	users, err := s.user.FindUser(ctx, roleId, nil, nil)
	if err != nil {
		return nil, helper.NewHTTPError(500, fmt.Errorf("failed to find users: %w", err))
	}

	if len(users) == 0 {
		return nil, helper.NewHTTPError(404, fmt.Errorf("users not found"))
	}

	var userResponses []*api.UserSchema
	for _, user := range users {
		userRoles := s.mapUserToUserRole(&user, nil)
		userResponses = append(userResponses, s.mapUserToUserSchema(userRoles))
	}

	return userResponses, nil
}

func (s *UserService) mapUserToUserSchema(user *model.UserRoles) *api.UserSchema {
	return &api.UserSchema{
		Id:    int(user.ID),
		Name:  user.Name,
		Email: user.Email,
		Roles: s.mapRolesToUserRole(user.Roles),
	}
}

func (s *UserService) mapRolesToUserRole(roles []model.Role) ( userRoles []api.RoleSchema) {
	for _, role := range roles {
		userRoles = append(userRoles, api.RoleSchema{
			Id:   int(role.ID),
			Name: role.Name,
		})
	}

	return 
}

func (s *UserService) mapUserToUserRole(user *model.User, roles []model.Role) *model.UserRoles {
	return &model.UserRoles{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Roles:     roles,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
