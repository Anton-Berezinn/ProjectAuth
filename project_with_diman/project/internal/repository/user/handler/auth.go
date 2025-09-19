package handler

import (
	"Project/internal/manager"
	"Project/internal/models"
	"Project/internal/repository/user/postgres"
	j "Project/internal/token/jwt"
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

var (
	ErrNotValid       = errors.New("not valid data")
	ErrDuplicateEmail = errors.New("duplicate email")
	ErrLenPassword    = errors.New("password is too short")
	ErrNotValidEmail  = errors.New("email is not valid")
	ErrNotValidName   = errors.New("last name and first name are required")
	re                = regexp.MustCompile(`^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$`)
)

func RegisterHandler(ctx context.Context, req *models.RegisterRequest, project *manager.Project) (string, error) {
	if len(req.Password) < 6 {
		return "", fmt.Errorf("Password is too short %w", ErrLenPassword)
	}
	if !re.MatchString(req.Email) {
		return "", fmt.Errorf("Email is not valid %w", ErrNotValidEmail)
	}
	if req.LastName == "" || req.FirstName == "" {
		return "", fmt.Errorf("LastName and FirstName are required %w", ErrNotValidName)
	}

	userId, err := project.UserService.CreateUser(ctx, req.FirstName, req.LastName, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, postgres.ErrDuplicateEntry) {
			return "", ErrDuplicateEmail
		}
		return "", fmt.Errorf("failed to create user: %w", err)
	}

	version := 1

	err = project.Cache.SetData(ctx, strconv.Itoa(userId), strconv.Itoa(version))
	if err != nil {
		project.Logger.Println("failed to set data from redis server", err)
	}

	token, err := j.NewJwt(ctx, userId, &version)
	if err != nil {
		project.Logger.Fatalln("Error to create token %s", err.Error())
		return "", err
	}

	return token, nil
}

func LoginHandler(ctx context.Context, req *models.LoginRequest, project *manager.Project) (string, error) {

	if len(req.Password) < 6 {
		return "", fmt.Errorf("Password is too short %w", ErrLenPassword)
	}

	if !re.MatchString(req.Email) {
		return "", fmt.Errorf("Email is not valid %w", ErrNotValidEmail)
	}

	userId, err := project.UserService.LoginUser(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, postgres.ErrNoRows) {
			return "", fmt.Errorf("user not found %w", ErrNotValid)
		}
		return "", fmt.Errorf("failed to login: %w", err)
	}
	data, err := project.Cache.GetData(ctx, strconv.Itoa(userId))
	version, err := strconv.Atoi(string(data))
	if err != nil {
		project.Logger.Fatalln("Error to Parse Value version from cache", err.Error())
		return "", err
	}

	userid := strconv.Itoa(userId)
	version += 1

	err = project.Cache.UpdateData(ctx, userid, version)

	token, err := j.NewJwt(ctx, userId, &version)
	if err != nil {
		project.Logger.Fatalln("Error to create token %s", err.Error())
		return "", err
	}
	return token, nil
}
