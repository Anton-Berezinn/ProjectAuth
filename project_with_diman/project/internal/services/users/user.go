package users

import (
	"Project/internal/configs"
	"Project/internal/repository/user/postgres"
	"context"
	"database/sql"
)

type UserService struct {
	db *postgres.DB
}

func NewUserService(db *postgres.DB) *UserService {
	return &UserService{db: db}
}

func NewDb(config *configs.Config) (*sql.DB, error) {
	return postgres.NewConnect(config)
}

func (s *UserService) CreateUser(ctx context.Context, firstname, lastname, email, password string) (int, error) {
	return s.db.CreateUser(ctx, firstname, lastname, email, password)
}

func (s *UserService) LoginUser(ctx context.Context, email, password string) (int, error) {
	return s.db.LoginUser(ctx, email, password)

}
