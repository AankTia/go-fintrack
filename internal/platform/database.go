package platform

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB adalah variabel global untuk menyimpan pool koneksi
// Kita gunakan variabel global agar bisa diakses oleh repository di seluruh aplikasi
var DB *pgxpool.Pool

// Connect menginisialisasi koneksi ke database
// Kita mengembalikan error jika koneksi gagal agar aplikasi bisa berhenti (fail-fast)
func Connect() error {
	// Ambil URL database dari environment variable (diatur di .env)
	connStr := os.Getenv("DB_URL")
	if connStr == "" {
		return fmt.Errorf("DB_URL is not se in environment")
	}

	// Buat pool
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Simpan ke variabel global
	DB = pool
	fmt.Println("Successfully connect to the database")
	return nil
}

// Close menutup semua koneksi dalam pool saat aplikasi dimatikan
func Close() {
	if DB != nil {
		DB.Close()
	}
}

// Ping memastikan koneksi ke database masih aktif
func Ping() error {
	ctx := context.Background()
	err := DB.Ping(ctx)
	if err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

func RunMigrations() error {
	// Ambil URL database dari environment variable (diatur di .env)
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		return fmt.Errorf("DB_URL is not se in environment")
	}

	m, err := migrate.New(
		"file://migrations",
		dbUrl,
	)
	if err != nil {
		return fmt.Errorf("migration init: %w", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("migration up: %w", err)
	}

	log.Println("✅ Database migrations applied")
	return nil
}
