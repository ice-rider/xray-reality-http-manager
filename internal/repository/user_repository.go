package repository

import (
	"database/sql"
	"errors"
	"xray_server/internal/domain"

	_ "github.com/mattn/go-sqlite3"
)

type UserRepositorySQLite struct {
	db *sql.DB
}

func NewUserRepositorySQLite(dbPath string) (*UserRepositorySQLite, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	repo := &UserRepositorySQLite{db: db}
	if err := repo.initDB(); err != nil {
		return nil, err
	}

	return repo, nil
}

func (r *UserRepositorySQLite) initDB() error {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL,
		role TEXT NOT NULL DEFAULT 'admin'
	);
	CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
	`
	_, err := r.db.Exec(query)
	return err
}

func (r *UserRepositorySQLite) FindByUsername(username string) (*domain.User, error) {
	query := `SELECT id, username, password, role FROM users WHERE username = ?`
	row := r.db.QueryRow(query, username)

	user := &domain.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Password, &user.Role)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepositorySQLite) Create(user *domain.User) error {
	query := `INSERT INTO users (username, password, role) VALUES (?, ?, ?)`
	result, err := r.db.Exec(query, user.Username, user.Password, user.Role)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = id
	return nil
}

func (r *UserRepositorySQLite) Close() error {
	return r.db.Close()
}
