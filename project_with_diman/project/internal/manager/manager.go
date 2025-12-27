package manager

import (
	"Project/internal/cache/redis"
	"Project/internal/repository/user/postgres"
	"Project/internal/services"
	"database/sql"
	"fmt"
	"log"
)

type Project struct {
	Db           *sql.DB
	Logger       *log.Logger
	UserService  *services.UserService
	CacheUser    *redis.RedisCache
	CacheCatalog *redis.RedisCache
}

func NewProject(db *sql.DB, logger *log.Logger, cacheUser *redis.RedisCache, cacheCatalog *redis.RedisCache) (*Project, error) {
	if db == nil || logger == nil {
		return nil, fmt.Errorf(`db or logger is nil`)
	}

	userRepo := postgres.NewDB(db)

	// Создаем сервис пользователей
	userService := services.NewUserService(userRepo)

	project := &Project{
		Db:           db,
		Logger:       logger,
		UserService:  userService,
		CacheUser:    cacheUser,
		CacheCatalog: cacheCatalog,
	}
	return project, nil
}

func (m *Project) SomeWork() {
	//Todo: Какая нибудь внутренняя работа
}
