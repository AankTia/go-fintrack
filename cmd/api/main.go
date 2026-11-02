package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AankTia/go-fintrack/internal/handler"
	"github.com/AankTia/go-fintrack/internal/platform"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cfg, err := platform.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// Connect ke Database
	if err := platform.Connect(); err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}

	// Verifikasi koneksi dengan Ping
	if err := platform.Ping(); err != nil {
		log.Fatalf("Database not ready: %v", err)
	}

	fmt.Println("Database connection verified!")

	if err := platform.RunMigrations(); err != nil {
		log.Fatalf("Could not run migrations: %v", err)
	}

	r := chi.NewRouter()

	// Tambahkan middleware “Logger” (meng‑log setiap request)
	r.Use(middleware.Logger)

	r.Use(middleware.Recoverer) // melaporkan panic menjadi 500
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello !"))
	})

	r.Route("/api/v1", func(r chi.Router){
		r.Post("/auth/register", handler.RegisterHandler)
		r.Post("/auth/login", handler.LoginHandler)
	})

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
