package service

import (
	"context"
	"errors"

	"github.com/AankTia/go-fintrack/internal/model"
	"github.com/AankTia/go-fintrack/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser mem‑proses pendaftaran
func RegisterUser(ctx context.Context, req *model.RegisterRequest) (*model.UserResponse, error) {
	if len(req.Email) == 0 {
		return nil, errors.New("email is required")
	}

	if len(req.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}
	// 1. Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// 2. Simpan ke DB (repo layer)
	user, err := repository.CreateUser(ctx, req.Email, string(hashed))
	if err != nil {
		// Check apakah error dari database karena email duplicate?
		// PostgreSQL: unique_violation = "23505"
		// if errors.Is(err, psqlErrorDuplicateEmail) {
		// 	return nil, errors.New("email already exists") // Return 400 Duplicate
		// }
		return nil, err
	}

	return &model.UserResponse{
		ID:    user.ID,
		Email: user.Email,
	}, nil
}

// LoginUser mem‑verifikasi login
func LoginUser(ctx context.Context, req *model.LoginRequest) (string, error) {
	user, err := repository.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	// TODO: generate JWT token (gunakan jwt-go atau libs)
	// Sebagai contoh sederhana:
	token := "mocked-jwt-token-placeholder"

	return token, nil
}
