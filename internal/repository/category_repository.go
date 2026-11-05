package repository

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/AankTia/go-fintrack/internal/model"
)

type CategoryRepository interface {
	Create(ctx context.Context, category *model.Category) error

	GetByUserId(ctx context.Context, userID string) ([]*model.Category, error)

	Delete(ctx context.Context, id string, userID string) error
}

type categoryRepository struct {
	db  *pgxpool.Pool
}

func NewCategoryRepository(db  *pgxpool.Pool) CategoryRepository {
	return &categoryRepository{
		db: db,
	}
}

func (r *categoryRepository) Create(ctx context.Context, category *model.Category) error {
	// Gunakan query parameterized ($1, $2, dst) untuk mencegah SQL Injection
	query := `INSERT INTO categories (user_id, name, type, icon)
						VALUE ($1, $2, $3, $4)
						RETURNING id`

	// ExecQuery digunakan karena kita butuh return value (id)
	err := r.db.QueryRow(ctx, query,
		category.UserID,
		category.Name,
		category.Type,
		category.Icon,
	).Scan(&category.ID) // Mengambil ID yang baru dibuat dari DB dan memasukkannya ke struct category

	return err
}

func (r *categoryRepository) GetByUserId(ctx context.Context, userID string) ([]*model.Category, error) {
	query := `SELECT id, user_id, name, type, icon, created_at, updated_at 
	          FROM categories WHERE user_id = $1 OR user_id IS NULL`
	
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close() // WAJIB tutup rows untuk mencegah memory leak

	var categories []*model.Category
	for rows.Next() {
		var cat model.Category
		if err := rows.Scan(&cat.ID, &cat.UserID, &cat.Name, &cat.Type, &cat.Icon, &cat.CreatedAt, &cat.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, &cat)
	}

	return categories, nil
}

func (r *categoryRepository) Delete(ctx context.Context, id string, userID string) error {
	query := `DELETE FROM categories WHERE id = $1 AND user_id = $2`

	rows, err := r.db.Query(ctx, query, id, userID)
	if err != nil {
		return err
	}
	defer rows.Close() // WAJIB tutup rows untuk mencegah memory leak

	return err
}
