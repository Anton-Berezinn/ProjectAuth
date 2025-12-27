package postgres

import (
	"Project/internal/configs"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
	"log"
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

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", fmt.Errorf("error to hash password: %w", err)
	}
	return string(bytes), err
}

func CompareHashAndPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

func NewConnect(config *configs.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=localhost port=5432 user=%s password=%s  sslmode=disable", config.Project.Db.Login, config.Project.Db.Password))
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

	pass, err := HashPassword(password)
	if err != nil {
		fmt.Println(err)
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}
	fmt.Println(pass)
	query := `INSERT INTO users (firstname, lastname, email, password) VALUES ($1, $2, $3, $4) RETURNING id`
	var id int
	err = tx.QueryRow(query, firstname, lastname, email, pass).Scan(&id)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" {
			return 0, fmt.Errorf("Duplicate email error: %w", ErrDuplicateEntry)
		}
		log.Println("failed to insert user:", err)
		return 0, fmt.Errorf("failed to create user: %w", err)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return id, nil
}

func (db *DB) GetUser(ctx context.Context, email, password string) (int, error) {
	tx, err := db.sql.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	row := db.sql.QueryRow("SELECT * FROM users WHERE email = $1", email)

	u := &User{}

	err = row.Scan(&u.Id, &u.LastName, &u.FirstName, &u.Email, &u.Password)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {

			return 0, fmt.Errorf("email or password is incorrect: %w", ErrNoRows)
		}
		log.Println("Error scan rows:", err)
		return 0, fmt.Errorf("failed to scan user: %w", err)
	}

	err = CompareHashAndPassword(u.Password, password)
	if err != nil {
		fmt.Println(err)
		return 0, fmt.Errorf("email or password is incorrect: %w", ErrNoRows)
	}

	if err = tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return u.Id, nil

}
