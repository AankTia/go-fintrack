package service

import (
	"context"
	"errors"
	"time"

	"github.com/AankTia/go-fintrack/internal/model"
	"github.com/AankTia/go-fintrack/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser mem‑proses pendaftaran
func RegisterUser(ctx context.Context, req *model.RegisterRequest) (*model.UserResponse, error) {
	// 1. Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// 2. Simpan ke DB (repo layer)
	user, err := repository.CreateUser(ctx, req.Email, string(hashed))
	if err != nil {
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
	token := "mocked-jwt-token"

	// Biasanya token ber‑ttl (exp) dan signed dengan secret
	_ = time.Now()

	return token, nil
}
