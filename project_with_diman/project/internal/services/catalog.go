package services

import (
	"Project/internal/models"
	"context"
)

func (u *UserService) GetCatalogService(ctx context.Context, filter string) (*[]models.CatalogResponse, error) {
	return u.db.GetCatalog(ctx, filter)
}

func (u *UserService) CreateCatalogService(ctx context.Context, data models.CreateCatalogRequest, userId int) (bool, error) {
	return u.db.CreateCatalog(ctx, data, userId)
}

func (u *UserService) UpdateCatalogService(ctx context.Context, data models.UpdateCatalogRequest, userId int) error {
	return u.db.UpdateCatalog(ctx, data, userId)

}
