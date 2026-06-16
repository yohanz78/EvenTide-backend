package delivery

import (
	"database/sql"
	"net/http"

	// Import layer lain dari domain Users
	"EventTide-backend/internal/users/repository"
	"EventTide-backend/internal/users/usecase"
)

// Fungsi ini yang akan dipanggil oleh main.go
func RegisterRoutes(mux *http.ServeMux, db *sql.DB) {
	// 1. Inisialisasi Repository
	userRepo := repository.NewPgUserRepository(db)

	// 2. Inisialisasi Usecase
	userUC := usecase.NewUserUsecase(userRepo)

	// 3. Inisialisasi Handler
	handler := &UserHandler{
		UserUC: userUC,
	}

	// 4. Daftarkan semua endpoint HTTP yang berkaitan dengan Users di sini
	mux.HandleFunc("POST /api/users/register", handler.RegisterUser)
	mux.HandleFunc("POST /api/users/login", handler.LoginUser)

	// Jika nanti ada endpoint login atau get profile, cukup tambahkan di sini:
	// mux.HandleFunc("POST /api/users/login", handler.LoginUser)
	// mux.HandleFunc("GET /api/users/profile", handler.GetProfile)
}
