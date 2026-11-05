package service

import (
	"context"
	"fmt"

	"github.com/AankTia/go-fintrack/internal/model"
	"github.com/AankTia/go-fintrack/internal/repository"
)

type CategoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) *CategoryService {
	return &CategoryService{
		repo: repo,
	}
}
func (s *CategoryService) CreateCategory(ctx context.Context, userID string, req model.CategoryRequest) (*model.Category, error) {
	// Business Validation
	if req.Name == "" {
		return nil, fmt.Errorf("%w: category name is required", ErrInvalidInput)
	}

	uID := userID
	category := &model.Category{
		UserID: &uID,
		Name:   req.Name,
		Type:   req.Type,
		Icon:   req.Icon,
	}

	// Panggil Repository untuk simpan ke DB
	err := s.repo.Create(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %v", err)
	}

	return category, nil
}

func (s *CategoryService) GetMyCategories(ctx context.Context, userID string) ([]*model.Category, error) {
	categories, err := s.repo.GetByUserId(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %v", err)
	}

	return categories, nil
}

func (s *CategoryService) DeleteCategory(ctx context.Context, id string, userID string) ([]*model.Category, error) {
	err := s.repo.Delete(ctx, id, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to delete category: %v", err)
	}

	return s.GetMyCategories(ctx, userID)
}
