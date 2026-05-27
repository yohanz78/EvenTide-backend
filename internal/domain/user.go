package domain

import (
	"context"
	"time"
)

// 1. STRUCT (Cetakan Data)
type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Tanda "-" memastikan password tidak akan pernah ikut terkirim ke output JSON
	CreatedAt    time.Time `json:"created_at"`
}

// INTERFACE

// 2. Kontrak Repository (Tugas Database)
// Usecase akan menggunakan ini tanpa perlu tahu apakah pakai PostgreSQL atau MySQL
type UserRepository interface {
	Save(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
}

// 3. Kontrak Usecase (Logika Bisnis Utama)
// Handler (Delivery) akan memanggil kontrak ini
type UserUsecase interface {
	Register(ctx context.Context, username, email, password string) (*User, error)
}
