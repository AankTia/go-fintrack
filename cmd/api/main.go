package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/AankTia/go-fintrack/internal/handler"
	"github.com/AankTia/go-fintrack/internal/middleware"
	"github.com/AankTia/go-fintrack/internal/repository"
	"github.com/AankTia/go-fintrack/internal/platform"
	"github.com/AankTia/go-fintrack/internal/service"
	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
)

func main() {
	// ---------------------------------------------------------
	// Load env Config
	// ---------------------------------------------------------
	cfg, err := platform.LoadConfig()
	if err != nil {
		log.Fatalf("Could not load config: %v", err)
	}

	// ---------------------------------------------------------
	// INISIALISASI DATABASE
	// ---------------------------------------------------------
	if err := platform.Connect(); err != nil {
		log.Fatalf("Could not connect to DB: %v", err)
	}
	defer platform.Close()

	// Verifikasi koneksi dengan Ping
	if err := platform.Ping(); err != nil {
		log.Fatalf("Database not ready: %v", err)
	}

	fmt.Println("Database connection verified!")

	if err := platform.RunMigrations(); err != nil {
		log.Fatalf("Could not run migrations: %v", err)
	}

	// ---------------------------------------------------------
	// DEPENDENCY INJECTION (THE WIRING)
	// ---------------------------------------------------------
	accountRepo := repository.NewAccountRepository(platform.DB)

	accountSvc := service.NewAccountService(accountRepo)

	accountHdl := handler.NewAccountHandler(accountSvc)

	// ---------------------------------------------------------
	// ROUTING
	// ---------------------------------------------------------
	r := chi.NewRouter()

	// Middleware Global
	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Recoverer) // melaporkan panic menjadi 500
	r.Use(chiMiddleware.AllowContentType("application/json"))
	r.Use(chiMiddleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello !"))
	})

	// Grup API v1
	r.Route("/api/v1", func(r chi.Router){
		// Route Public 
		r.Post("/auth/register", handler.RegisterHandler)
		r.Post("/auth/login", handler.LoginHandler)

		// Route Private (Wajib Login)
		r.Group(func(r chi.Router) {
			r.Use(middleware.AuthMiddleware) // Middleware yang menyuntikkan user_id ke context

			r.Post("/accounts", accountHdl.Create)
			r.Get("/accounts", accountHdl.List)
			r.Get("/accounts/{id}", accountHdl.GetDetail)
		})
	})

	addr := fmt.Sprintf(":%s", cfg.Port)
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
