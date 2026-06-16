package user

import (
	"database/sql"
	"errors"
	"log"
	"time"
)

type postgresRepo struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &postgresRepo{db: db}
}

func Migrate(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id         SERIAL PRIMARY KEY,
		name       TEXT        NOT NULL,
		email      TEXT        NOT NULL UNIQUE,
		is_active  BOOLEAN     NOT NULL DEFAULT true,
		created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);`)
	return err
}

func (r *postgresRepo) Create(u *User) error {
	now := time.Now()
	return r.db.QueryRow(
		`INSERT INTO users (name, email, is_active, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5)
		 RETURNING id, created_at, updated_at`,
		u.Name, u.Email, u.IsActive, now, now,
	).Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
}

func (r *postgresRepo) FindAll() ([]User, error) {
	rows, err := r.db.Query(
		`SELECT id, name, email, is_active, created_at, updated_at FROM users`)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("close rows: %v", err)
		}
	}()
	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *postgresRepo) FindByID(id int) (*User, error) {
	var u User
	err := r.db.QueryRow(
		`SELECT id, name, email, is_active, created_at, updated_at FROM users WHERE id=$1`, id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.IsActive, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, errors.New("user not found")
	}
	return &u, err
}

func (r *postgresRepo) FindByActive(active bool) ([]User, error) {
	rows, err := r.db.Query(
		`SELECT id, name, email, is_active, created_at, updated_at FROM users WHERE is_active=$1`, active)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("close rows: %v", err)
		}
	}()
	var users []User
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.IsActive, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *postgresRepo) Count() (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&count)
	return count, err
}

func (r *postgresRepo) Update(u *User) error {
	u.UpdatedAt = time.Now()
	result, err := r.db.Exec(
		`UPDATE users SET name=$1, email=$2, is_active=$3, updated_at=$4 WHERE id=$5`,
		u.Name, u.Email, u.IsActive, u.UpdatedAt, u.ID,
	)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("user not found")
	}
	return nil
}

func (r *postgresRepo) Delete(id int) error {
	result, err := r.db.Exec(`DELETE FROM users WHERE id=$1`, id)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return errors.New("user not found")
	}
	return nil
}
