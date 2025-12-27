package handler

import (
	"Project/internal/cache/redis"
	"Project/internal/manager"
	"Project/internal/models"
	"Project/internal/repository/user/postgres"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"
)

var (
	ErrClient        = errors.New("failed to update data")
	ErrAlreadyExists = errors.New("already exists")
)

func Validate(ctx context.Context, filter string) (string, bool) {

	value := strings.ToLower(strings.Fields(filter)[0])

	if value != "одежда" && value != "инструменты" && value != "товары для дома" {
		log.Println("")
		return "", false
	}
	return value, true
}

func GetBucketHandler(ctx context.Context, project *manager.Project, userID int) (*models.CreateBuckerUserId, error) {
	userId := strconv.Itoa(userID)
	value, err := project.CacheCatalog.GetCatalogs(ctx, userId)
	if err != nil {
		if errors.Is(err, redis.ErrNotFound) {
			user := &models.CreateBuckerUserId{
				UserId:    userID,
				Id_bucket: make([]int, 0, 100),
				Count:     0,
			}
			value, err := json.Marshal(user)
			if err != nil {
				project.Logger.Println("Error to marshal user %s", err)
				return nil, err
			}
			err = project.CacheCatalog.InsertCatalog(ctx, string(value), userId)
			if err != nil {
				project.Logger.Println("Error to insert catalog %s", err)
				return nil, err
			}
			return nil, nil
		}
		project.Logger.Println("Error unknow From GetCatalogs %s", err)
		return nil, err
	}
	user := &models.CreateBuckerUserId{}
	if err := json.Unmarshal([]byte(value), &user); err != nil {
		project.Logger.Println("Error to unmarshal user %s", err)
		return nil, err
	}
	fmt.Println(user)
	return user, nil
}

func CreateBucketHandler(ctx context.Context, project *manager.Project, id, userID int) error {
	userId := strconv.Itoa(userID)
	value, err := project.CacheCatalog.GetCatalogs(ctx, userId)
	if err != nil {
		if errors.Is(err, redis.ErrNotFound) {
			user := &models.CreateBuckerUserId{
				UserId:    userID,
				Id_bucket: make([]int, 0, 100),
				Count:     1,
			}
			user.Id_bucket = append(user.Id_bucket, id)
			value, err := json.Marshal(user)
			if err != nil {
				project.Logger.Println("Error to marshal user %s", err)
				return err
			}
			err = project.CacheCatalog.InsertCatalog(ctx, string(value), userId)
			if err != nil {
				project.Logger.Println("Error to insert catalog %s", err)
				return err
			}
			return nil
		}
		project.Logger.Println("Error unknow From GetCatalogs %s", err)
		return err
	}
	user := &models.CreateBuckerUserId{}
	if err := json.Unmarshal([]byte(value), &user); err != nil {
		project.Logger.Println("Error to unmarshal user %s", err)
		return err
	}
	ok := slices.Contains(user.Id_bucket, id)
	if ok {
		return fmt.Errorf("Already exists %w", ErrAlreadyExists)
	}

	user.Id_bucket = append(user.Id_bucket, id)
	user.Count++
	data, err := json.Marshal(user)
	if err != nil {
		project.Logger.Println("Error to marshal user %s", err)
		return err
	}
	err = project.CacheCatalog.InsertCatalog(ctx, string(data), userId)
	if err != nil {
		project.Logger.Println("Error to insert catalog %s", err)
		return err
	}
	return nil
}

func CreateCatalogHandler(ctx context.Context, data models.CreateCatalogRequest, userId int, project *manager.Project) (bool, error) {
	ok, err := project.UserService.CreateCatalogService(ctx, data, userId)
	if err != nil {
		project.Logger.Println("Error CreateCatalogService error: ", err)
		return false, err
	}
	if !ok {
		return false, errors.New("<UNK>")
	}
	return true, nil
}

func UpdateCatalogHandler(ctx context.Context, data models.UpdateCatalogRequest, userId int, project *manager.Project) error {
	err := project.UserService.UpdateCatalogService(ctx, data, userId)
	if err != nil {
		if errors.Is(err, postgres.ErrClient) {
			return fmt.Errorf("%w", ErrClient)
		}
		project.Logger.Println("Error UpdateCatalogService error: ", err)
		return err
	}
	return nil
}
func GetCatalogHandler(ctx context.Context, filter string, project *manager.Project) (*[]models.CatalogResponse, error) {
	value, err := project.UserService.GetCatalogService(ctx, filter)
	if err != nil {
		project.Logger.Println("GetCatalogHandler err:", err)
		return nil, err
	}
	return value, nil
}
