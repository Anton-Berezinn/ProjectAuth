package manager

import (
	"Project/internal/repository/user/postgres"
	"Project/internal/services/users"
	"database/sql"
	"fmt"
	"log"
)

type Project struct {
	Db          *sql.DB
	Logger      *log.Logger
	UserService *users.UserService
}

func NewProject(db *sql.DB, logger *log.Logger) (*Project, error) {
	if db == nil || logger == nil {
		return nil, fmt.Errorf(`db or logger is nil`)
	}

	userRepo := postgres.NewDB(db)

	// Создаем сервис пользователей
	userService := users.NewUserService(userRepo)

	project := &Project{
		Db:          db,
		Logger:      logger,
		UserService: userService,
	}
	return project, nil
}

func (m *Project) SomeWork() {
	//Todo: Какая нибудь внутренняя работа
}
