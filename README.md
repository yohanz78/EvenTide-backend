# EventTide - Event Management Platform

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-4169E1?style=for-the-badge&logo=postgresql&logoColor=white)

## Overview

EvenTide is a service for blogging, event management, and documentation platform. It allows users to create and manage events, upload media, and share content seamlessly. The backend is built using Go, with PostgreSQL as the database.

## Project Structure

```
eventide-backend/
├── cmd/
│   └── api/
│       └── main.go                 // Titik masuk. Hanya berisi wiring DB, Server, dan Router.
├── internal/
│   ├── domain/                     // (Layer 1) Kumpulan Structs dan Interfaces untuk semua entitas
│   │   ├── user.go
│   │   ├── event.go
│   │   └── media.go
│   ├── events/                     // Domain Acara
│   │   ├── delivery/
│   │   │   └── http_handler.go     // (Layer 4) Menerima request HTTP, membalas JSON
│   │   ├── usecase/
│   │   │   └── event_usecase.go    // (Layer 3) Logika bisnis (misal: validasi tanggal acara)
│   │   └── repository/
│   │       └── pg_repository.go    // (Layer 2) Query SQL ke PostgreSQL khusus untuk Event
│   ├── media/                      // Domain Dokumentasi
│   │   ├── delivery/
│   │   │   └── http_handler.go
│   │   ├── usecase/
│   │   │   └── media_usecase.go    // Logika bisnis (misal: memanggil API Cloud Storage)
│   │   └── repository/
│   │       └── pg_repository.go    // Query SQL untuk menyimpan URL media
│   └── users/                      // Domain Pengguna (Auth)
│       └── ... (struktur serupa)
├── pkg/                            // Library bantuan yang bisa dipakai di mana saja
│   ├── database/                   // Setup koneksi PostgreSQL
│   │   └── postgres.go
│   └── storage/                    // Konfigurasi koneksi ke Cloud Storage (AWS/GCP)
│       └── s3_client.go
├── db/
│   └── migrations/                 // Tempat menyimpan file-file SQL
├── .env                            // Variabel lingkungan (Kredensial)
└── go.mod
```

## Commands

- **Run Migrations**
    - Up:

        ```bash
        migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/your_db?sslmode=disable" -verbose up
        ```

        or

        ```bash
        go run cmd/migrate/main.go
        ```

    - Down:

        ```bash
        migrate -path db/migrations -database "postgresql://postgres:password@localhost:5432/your_db?sslmode=disable" -verbose down
        ```

- **Run Server**:

    ```bash
    go run cmd/api/main.go
    ```

    or

    ```bash
    air
    ```
