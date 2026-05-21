package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"EventTide-backend/pkg/database"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
)

// 1. Fungsi Migrasi
func runDBMigration(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("❌ Gagal inisialisasi driver migrasi: %v", err)
	}

	// Membaca file SQL dari folder db/migrations (sesuaikan path relatifnya)
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf("❌ Gagal membaca folder migrasi: %v", err)
	}

	// Eksekusi migrasi Up
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("❌ Gagal mengeksekusi migrasi: %v", err)
	}

	fmt.Println("✅ [Migration] Skema database berhasil diperbarui (atau sudah up-to-date)!")
}

func main() {
	fmt.Println("=== Memulai Eksekusi Database Migrations ===")

	// 1. Memuat konfigurasi dari file .env di root folder
	if err := godotenv.Load(".env"); err != nil {
		log.Println("⚠️ File .env tidak ditemukan, menggunakan environment sistem.")
	}

	// 2. Menggunakan fungsi ConnectDB dari package database kita
	db, err := database.ConnectDB(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)
	if err != nil {
		log.Fatalf("❌ Gagal terhubung ke database: %v", err)
	}
	defer db.Close() // Tutup koneksi saat migrasi selesai

	// 3. Konfigurasi driver spesifik untuk migrasi PostgreSQL
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("❌ Gagal inisialisasi driver migrasi: %v", err)
	}

	// 4. Mengarahkan path ke folder db/migrations (dibaca dari root proyek)
	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf("❌ Gagal membaca folder migrasi: %v", err)
	}

	// 5. Mengeksekusi perintah UP
	err = m.Up()

	// 6. Penanganan Error yang Rapi
	if err != nil {
		// ErrNoChange BUKAN error. Itu artinya DB sudah dalam versi terbaru.
		if err == migrate.ErrNoChange {
			fmt.Println("✅ Skema database sudah up-to-date. Tidak ada eksekusi baru.")
		} else {
			// Jika error lain (seperti salah sintaks SQL), program berhenti
			log.Fatalf("❌ Gagal mengeksekusi migrasi: %v", err)
		}
	} else {
		fmt.Println("✅ Migrasi sukses! Skema database berhasil diperbarui.")
	}

	fmt.Println("=== Proses Migrasi Selesai ===")
}
