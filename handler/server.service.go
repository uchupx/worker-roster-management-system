package handler

import (
	"github.com/uchupx/worker-roster-management-system/internal/service"
	"github.com/uchupx/worker-roster-management-system/internal/service/jwt"
)

func (s *Server) getJWTService() *jwt.CryptService {
	if s.jwtService != nil {
		return s.jwtService
	}

	svc := jwt.NewCryptService(jwt.Params{
		Conf: jwt.JWTConf{
			Privatekey: s.conf.RSA.Private,
			PublicKey:  s.conf.RSA.Public,
		},
	})

	s.jwtService = &svc

	return s.jwtService
}

func (s *Server) getAuthService() *service.AuthService {
	if s.authService != nil {
		return s.authService
	}

	s.authService = service.NewAuthService(service.AuthServiceParams{
		User:        s.getUserRepo(),
		JWT:         *s.getJWTService(),
		RedisClient: s.getRedisClient(),
	})

	return s.authService
}

func (s *Server) getUserService() *service.UserService {
	if s.userService != nil {
		return s.userService
	}

	s.userService = service.NewUserService(service.UserServiceParams{
		User:        s.getUserRepo(),
		JWT:         *s.getJWTService(),
		RoleService: *s.getRoleService(),
	})

	return s.userService
}

func (s *Server) getShiftService() *service.ShiftService {
	if s.shiftService != nil {
		return s.shiftService
	}

	s.shiftService = service.NewShiftService(service.ShiftServiceParams{
		UserService:     *s.getUserService(),
		ShiftRepository: s.getShiftRepo(),
		RoleService:     *s.getRoleService(),
	})

	return s.shiftService
}

func (s *Server) getRoleService() *service.RoleService {
	if s.roleService != nil {
		return s.roleService
	}

	s.roleService = service.NewRoleService(service.RoleServiceParams{
		RoleRepository: s.getRoleRepo(),
	})

	return s.roleService
}
