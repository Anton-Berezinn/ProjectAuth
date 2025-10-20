package postgres

import (
	"Project/internal/models"
	"context"
	"errors"
	"fmt"
	"log"
)

var (
	ErrServer = errors.New("server error")
	ErrClient = errors.New("client error")
)

func (db *DB) CreateCatalog(ctx context.Context, data models.CreateCatalogRequest, userId int) (bool, error) {
	ok, err := db.sql.Exec(`INSERT INTO products(name,description,filter,price,userid) VALUES ($1,$2,$3,$4,$5)`, data.Name, data.Description, data.Filter, data.Price, userId)
	if err != nil {
		return false, fmt.Errorf("failed to insert catalog: %w", ErrServer)
	}
	value, err := ok.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("failed to get rows affected: %w", err)
	}
	if value != 1 {
		return false, fmt.Errorf("expected 1 row affected, got %d", value)
	}
	return true, nil
}

func (db *DB) GetCatalog(ctx context.Context, filter string) (*[]models.CatalogResponse, error) {

	value, err := db.sql.Query(`SELECT id,name,description,price FROM products WHERE filter=$1 ORDER BY popularity DESC LIMIT 5`, filter)
	if err != nil {
		return nil, fmt.Errorf("get catalog: %w", err)
	}

	defer value.Close()

	values := make([]models.CatalogResponse, 0)
	for value.Next() {

		data := models.CatalogResponse{}
		if err := value.Scan(&data.Id, &data.Name, &data.Description, &data.Price); err != nil {
			return nil, fmt.Errorf("get catalog: %w", err)
		}
		values = append(values, data)
	}
	return &values, nil
}

func (db *DB) UpdateCatalog(ctx context.Context, data models.UpdateCatalogRequest, userId int) error {
	fmt.Println(data)
	val, err := db.sql.Exec("UPDATE products SET  name = $1, description = $2, price = $3, filter = $4 WHERE id = $5 and userid = $6", data.Name, data.Description, data.Price, data.Filter, data.Id, userId)
	if err != nil {
		log.Println("Failed to updateCatalog %s", err)
		return fmt.Errorf("failed to update catalog: %w", ErrServer)
	}

	value, err := val.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	if value != 1 {
		return fmt.Errorf("Failed to update catalog userId %w", ErrClient)
	}
	return nil
}
