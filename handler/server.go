package handler

import (
	"fmt"
	"log"
	"os"

	redis "github.com/redis/go-redis/v9"
	"github.com/uchupx/kajian-api/pkg/db"

	"github.com/uchupx/worker-roster-management-system/config"
	"github.com/uchupx/worker-roster-management-system/internal/database"
	"github.com/uchupx/worker-roster-management-system/internal/repository"
	"github.com/uchupx/worker-roster-management-system/internal/service"
	"github.com/uchupx/worker-roster-management-system/internal/service/jwt"
)

type Server struct {
	conf *config.Config

	db          *db.DB
	redisClient *redis.Client

	authService  *service.AuthService
	jwtService   *jwt.CryptService
	userService  *service.UserService
	shiftService *service.ShiftService
	roleService  *service.RoleService

	userRepo  repository.UserRepositoryInterface
	shiftRepo repository.ShiftRepositoryInterface
	roleRepo  repository.RoleRepositoryInterface
}

func NewServer(conf *config.Config) *Server {
	server := Server{
		conf: conf,
	}

	server.initDB()
	server.initSchema()

	server.getAuthService()
	server.getUserService()
	server.getShiftService()
	server.getRoleService()

	return &server
}

func (s *Server) initDB() *db.DB {
	if s.db != nil {
		return s.db
	}

	conn, err := database.NewConnection(database.DBPayload{
		Database: s.conf.App.Database,
	})

	if err != nil {
		panic(err)
	}

	s.db = conn
	fmt.Println("Connected to database:", s.conf.App.Database)
	return s.db
}

func (s Server) initSchema() {
	conn := s.initDB()

	sqlBytes, err := os.ReadFile(s.conf.App.Schema)
	if err != nil {
		log.Fatalf("Failed to read SQL file: %v", err)
	}

	sqlContent := string(sqlBytes)
	if _, err = conn.Exec(sqlContent); err != nil {
		panic(fmt.Errorf("failed to execute schema SQL file: %v", err))
	}
}

func (s *Server) getRedisClient() *redis.Client {
	if s.redisClient == nil {
		s.redisClient = database.GetRedisClient(*s.conf)
	}
	return s.redisClient
}
