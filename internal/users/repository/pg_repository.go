package repository

import (
	"context"
	"database/sql"
	"errors"

	// Import domain untuk mengambil Struct dan Interface
	"EventTide-backend/internal/domain"
)

type pgUserRepository struct {
	DB *sql.DB
}

// Fungsi "Pabrik" untuk membuat repository baru
func NewPgUserRepository(db *sql.DB) domain.UserRepository {
	return &pgUserRepository{DB: db}
}

func (r *pgUserRepository) Save(ctx context.Context, user *domain.User) error {
	query := `INSERT INTO users (username, email, password_hash) 
			  VALUES ($1, $2, $3) RETURNING id, created_at`

	// QueryRowContext memastikan eksekusi dibatalkan jika 'ctx' dibatalkan
	err := r.DB.QueryRowContext(ctx, query, user.Username, user.Email, user.PasswordHash).
		Scan(&user.ID, &user.CreatedAt)

	return err
}

func (r *pgUserRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `SELECT id, username, email, password_hash, created_at FROM users WHERE email = $1`

	user := &domain.User{}
	err := r.DB.QueryRowContext(ctx, query, email).
		Scan(&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.CreatedAt)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Email tidak ditemukan (bukan error sistem)
		}
		return nil, err // Error sistem/koneksi
	}
	return user, nil
}
