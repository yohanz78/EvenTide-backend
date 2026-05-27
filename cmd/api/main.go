package main

import (
	"EventTide-backend/pkg/database"
	"fmt"
	"log"
	"net/http"
	"os"

	usersDelivery "EventTide-backend/internal/users/delivery"

	"github.com/joho/godotenv"
)

// 2. Titik Masuk Utama
func main() {
	fmt.Println("=== Memulai Event & Media Platform API ===")

	// 1. Memuat konfigurasi dari file .env
	if err := godotenv.Load(".env"); err != nil {
		log.Println("⚠️ File .env tidak ditemukan, menggunakan environment sistem.")
	}

	// 2. Membuka koneksi database yang akan dipakai terus-menerus oleh API
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

	// 3. Setup Router & Server HTTP
	mux := http.NewServeMux()

	// 4. Mendaftarkan Route
	usersDelivery.RegisterRoutes(mux, db)

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
