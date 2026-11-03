package repository

import (
	"context"
	"fmt"

	"github.com/AankTia/go-fintrack/internal/platform" // untuk DB pool
)

// User struct (model internal)
type User struct {
	ID           int64
	Email        string
	PasswordHash string
}

// CreateUser men-insert user baru ke database dan return user yang tersimpan
// Parameter: email, hashedPassword (password sudah di-hash oleh service layer)
// Return: object User dengan ID baru, atau error jika gagal insert
func CreateUser(ctx context.Context, email, hashedPassword string) (*User, error) {
	// Query SQL dengan parameterized query (cegah SQL Injection)
	//    - Menggunakan $1, $2 untuk placeholder parameter
	//    - RETURNING clause ambil id & email dari DB sebagai primary response
	query := `INSERT INTO users (email, password_hash, created_at, updated_at)
						VALUES ($1, $2, NOW(), NOW())
						RETURNING id, email`

	var user User
	if err := platform.DB.QueryRow(ctx, query, email, hashedPassword).Scan(&user.ID, &user.Email); err != nil {
		return nil, fmt.Errorf("cannot create user: %w", err)
	}

	return &user, nil
}

// GetUserByEmail mencari user berdasar email
func GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, email, password_hash FROM users WHERE email=$1`

	var user User
	if err := platform.DB.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.PasswordHash); err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	return &user, nil
}
