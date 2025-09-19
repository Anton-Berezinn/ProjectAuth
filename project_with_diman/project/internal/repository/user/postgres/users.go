package postgres

import (
	"Project/internal/configs"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type DB struct {
	sql *sql.DB
}

var (
	ErrOpen           = errors.New("failed to open database")
	ErrDuplicateEntry = errors.New("duplicate entry")
	ErrNoRows         = errors.New("no rows ")
)

type User struct {
	Id        int
	LastName  string
	FirstName string
	Email     string
	Password  string
}

func NewConnect(config *configs.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=users sslmode=disable", config.Project.Db.Login, config.Project.Db.Password))
	if err != nil {
		return nil, fmt.Errorf(`failed to open database: %w`, ErrOpen)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf(`failed to ping database: %w`, err)
	}

	return db, nil
}

func NewDB(db *sql.DB) *DB {
	return &DB{sql: db}
}

func (db *DB) CreateUser(ctx context.Context, firstname, lastname, email, password string) (int, error) {
	tx, err := db.sql.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO users (firstname, lastname, email, password) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err = tx.QueryRow(query, firstname, lastname, email, password).Scan(&id)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			// 23505 — уникальное ограничение нарушено
			return 0, fmt.Errorf("Duplicate email error: %w", ErrDuplicateEntry)
		}

		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return id, nil
}

func (db *DB) LoginUser(ctx context.Context, email, password string) (int, error) {
	tx, err := db.sql.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	row := db.sql.QueryRow("SELECT * FROM users WHERE email = $1", email)

	u := &User{}

	err = row.Scan(&u.Id, &u.LastName, &u.FirstName, &u.Password, &u.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return 0, fmt.Errorf("email or password is incorrect: %w", ErrNoRows)
		}
		return 0, fmt.Errorf("failed to scan user: %w", err)
	}
	if u.Password != password {
		return 0, fmt.Errorf("email or password is incorrect: %w", ErrNoRows)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return u.Id, nil

}
