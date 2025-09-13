package handler

import (
	"Project/internal/manager"
	"Project/internal/models"
	"Project/internal/repository/user/postgres"
	"context"
	"errors"
	"fmt"
	"regexp"
)

var (
	ErrNotValid       = errors.New("not valid data")
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrLenPassword    = errors.New("password is too short")
	ErrNotValidEmail  = errors.New("email is not valid")
	ErrNotValidName   = errors.New("last name and first name are required")
	re                = regexp.MustCompile(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`)
)

func RegisterHandler(ctx context.Context, req *models.RegisterRequest, project *manager.Project) (int, error) {
	if len(req.Password) < 6 {
		return 0, fmt.Errorf("Password is too short %w", ErrLenPassword)
	}
	if !re.MatchString(req.Email) {
		return 0, fmt.Errorf("Email is not valid %w", ErrNotValidEmail)
	}
	if req.LastName == "" || req.FirstName == "" {
		return 0, fmt.Errorf("LastName and FirstName are required %w", ErrNotValidName)
	}

	userId, err := project.UserService.CreateUser(ctx, req.FirstName, req.LastName, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, postgres.ErrDuplicateEntry) {
			return 0, ErrDuplicateEmail
		}
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	return userId, nil
}

func LoginHandler(ctx context.Context, req *models.LoginRequest, project *manager.Project) (int, error) {

	if len(req.Password) < 6 {
		return 0, fmt.Errorf("Password is too short %w", ErrLenPassword)
	}

	if !re.MatchString(req.Email) {
		return 0, fmt.Errorf("Email is not valid %w", ErrNotValidEmail)
	}

	userId, err := project.UserService.LoginUser(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, postgres.ErrNoRows) {
			return 0, fmt.Errorf("user not found %w", ErrNotValid)
		}
		return 0, fmt.Errorf("failed to login: %w", err)
	}
	return userId, nil
}
