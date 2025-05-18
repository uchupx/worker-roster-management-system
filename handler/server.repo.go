package handler

import (
	"github.com/uchupx/worker-roster-management-system/internal/repository"
	"github.com/uchupx/worker-roster-management-system/internal/repository/sqlite"
)

func (s *Server) getUserRepo() repository.UserRepositoryInterface {
	if s.userRepo != nil {
		return s.userRepo
	}

	s.userRepo = sqlite.NewUserRepository(s.initDB())

	return s.userRepo
}

func (s *Server) getShiftRepo() repository.ShiftRepositoryInterface {
	if s.shiftRepo != nil {
		return s.shiftRepo
	}

	s.shiftRepo = sqlite.NewShiftRepository(s.initDB())

	return s.shiftRepo
}

func (s *Server) getRoleRepo() repository.RoleRepositoryInterface {
	if s.roleRepo != nil {
		return s.roleRepo
	}

	s.roleRepo = sqlite.NewRoleRepository(s.initDB())

	return s.roleRepo
}
