package database

import (
	"database/sql"
	"fmt"

	// Import driver PostgreSQL (wajib ada di file yang memanggil sql.Open)
	_ "github.com/lib/pq"
)

// Fungsi ConnectDB menerima kredensial, membuka koneksi, dan mengembalikan objek *sql.DB
func ConnectDB(host, port, user, password, dbName, sslMode string) (*sql.DB, error) {
	// Merakit Connection String
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbName, sslMode)

	// Membuka koneksi
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		// fmt.Errorf untuk memberikan konteks error yang jelas
		return nil, fmt.Errorf("Gagal inisialisasi koneksi: %w", err)
	}

	// Memastikan server database benar-benar menyala dan merespons
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("Database tidak merespons (ping gagal): %w", err)
	}

	fmt.Println("✅ Koneksi fisik ke PostgreSQL berhasil!")

	// Mengembalikan objek db agar bisa digunakan oleh main.go
	return db, nil
}
