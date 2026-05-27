package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"EventTide-backend/internal/users/delivery"
	"EventTide-backend/internal/users/repository"
	"EventTide-backend/internal/users/usecase"
	"EventTide-backend/pkg/database"

	"github.com/joho/godotenv"
)

// 2. Titik Masuk Utama
func main() {
	fmt.Println("=== Memulai Event & Media Platform API ===")

	// Memuat konfigurasi dari file .env
	if err := godotenv.Load(".env"); err != nil {
		log.Println("⚠️ File .env tidak ditemukan, menggunakan environment sistem.")
	}

	// Membuka koneksi database yang akan dipakai terus-menerus oleh API
	db, err := database.ConnectDB(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)
	if err != nil {
		log.Fatalf("❌ Aplikasi berhenti, gagal konek DB: %v", err)
	}
	defer db.Close() // Menjaga koneksi tetap hidup sampai server dimatikan

	// 1. Inisialisasi Repository
	userRepo := repository.NewPgUserRepository(db)

	// 2. Inisialisasi Usecase (menyuntikkan Repo)
	userUsecase := usecase.NewUserUsecase(userRepo)

	// 3. Inisialisasi Handler (menyuntikkan Usecase)
	userHandler := &delivery.UserHandler{
		UserUC: userUsecase,
	}

	// Setup Router & Server HTTP
	mux := http.NewServeMux()

	// 4. Mendaftarkan Route
	mux.HandleFunc("POST /api/users/register", userHandler.RegisterUser)

	// Endpoint percobaan
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status": "OK", "message": "API siap menerima request!"}`))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🚀 Server running: http://localhost:%s\n", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("❌ Gagal menjalankan server: %v", err)
	}
}
