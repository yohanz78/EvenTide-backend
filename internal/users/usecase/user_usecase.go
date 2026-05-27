package usecase

import (
	"context"
	"errors"

	"EventTide-backend/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo domain.UserRepository
}

// Inject Repository ke dalam Usecase
func NewUserUsecase(ur domain.UserRepository) domain.UserUsecase {
	return &userUsecase{userRepo: ur}
}

func (u *userUsecase) Register(ctx context.Context, username, email, password string) (*domain.User, error) {
	// 1. Validasi Bisnis: Cek apakah email sudah terdaftar
	existingUser, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("email sudah terdaftar")
	}

	// 2. Keamanan: Hash password agar tidak bocor jika database diretas
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 3. Merakit entitas baru
	newUser := &domain.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
	}

	// 4. Perintahkan Repository untuk menyimpan ke database
	if err := u.userRepo.Save(ctx, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}
