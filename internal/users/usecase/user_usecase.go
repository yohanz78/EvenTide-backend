package usecase

import (
	"context"
	"errors"
	"os"
	"time"

	"EventTide-backend/internal/domain"

	jwt "github.com/golang-jwt/jwt/v5"
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

func (u *userUsecase) Login(ctx context.Context, email, password string) (string, error) {
	// 1. Cari user berdasarkan email di database
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("email atau password salah") // Jangan pernah sebut spesifik "email tidak ditemukan" demi keamanan
	}

	// 2. Bandingkan password murni dengan password hash di database
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("email atau password salah")
	}

	// 3. Pembuatan Gelang VIP (JWT)
	// Claims adalah data apa saja yang mau kita "titipkan" di dalam token tersebut
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token kadaluarsa dalam 24 jam
	}

	// Membuat struktur token menggunakan algoritma HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 4. Stempel token tersebut menggunakan JWT_SECRET dari .env
	secretKey := []byte(os.Getenv("JWT_SECRET"))
	signedToken, err := token.SignedString(secretKey)
	if err != nil {
		return "", errors.New("gagal membuat token keamanan")
	}

	return signedToken, nil
}
