package service

import (
	"context"
	"testing"

	"github.com/AankTia/go-fintrack/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCategoryRepository adalah mock untuk CategoryRepository
type MockCategoryRepository struct {
	mock.Mock
}

func (m *MockCategoryRepository) Create(ctx context.Context, category *model.Category) error {
	args := m.Called(ctx, category)
	return args.Error(0)
}

func (m *MockCategoryRepository) GetByUserId(ctx context.Context, userID string) ([]*model.Category, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Category), args.Error(1)
}

func (m *MockCategoryRepository) Delete(ctx context.Context, id string, userID string) error {
	args := m.Called(ctx, id, userID)
	return args.Error(0)
}

func TestCreateCategory(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	svc := NewCategoryService(mockRepo)
	ctx := context.Background()
	userID := "user-123"

	t.Run("success create category", func(t *testing.T) {
		req := model.CategoryRequest{
			Name: "Makanan",
			Type: "expense",
			Icon: "food-icon",
		}

		// Mock: User belum punya kategori dengan nama "Makanan"
		mockRepo.On("GetByUserId", ctx, userID).Return([]*model.Category{}, nil).Once()
		mockRepo.On("Create", ctx, mock.AnythingOfType("*model.Category")).Return(nil).Once()

		cat, err := svc.CreateCategory(ctx, userID, req)

		assert.NoError(t, err)
		assert.Equal(t, "Makanan", cat.Name)
		mockRepo.AssertExpectations(t)
	})

	t.Run("fail duplicate category name", func(t *testing.T) {
		req := model.CategoryRequest{
			Name: "Makanan",
			Type: "expense",
		}

		// Mock: User SUDAH punya kategori bernama "Makanan"
		existingCats := []*model.Category{
			{ID: "cat-1", Name: "Makanan", Type: "expense"},
		}
		mockRepo.On("GetByUserId", ctx, userID).Return(existingCats, nil).Once()

		cat, err := svc.CreateCategory(ctx, userID, req)

		assert.Error(t, err)
		assert.Nil(t, cat)
		assert.Contains(t, err.Error(), "category name already exists for this user")
	})

	t.Run("fail empty name", func(t *testing.T) {
		req := model.CategoryRequest{Name: ""}

		cat, err := svc.CreateCategory(ctx, userID, req)

		assert.Error(t, err)
		assert.Nil(t, cat)
		assert.Contains(t, err.Error(), "category name is required")
	})
}

func TestDeleteCategory(t *testing.T) {
	mockRepo := new(MockCategoryRepository)
	svc := NewCategoryService(mockRepo)
	ctx := context.Background()
	userID := "user-123"

	t.Run("success delete category", func(t *testing.T) {
		catID := "cat-abc"

		// Mock: Berhasil hapus, lalu return list kategori terbaru
		mockRepo.On("Delete", ctx, catID, userID).Return(nil).Once()
		mockRepo.On("GetByUserId", ctx, userID).Return([]*model.Category{}, nil).Once()

		cats, err := svc.DeleteCategory(ctx, catID, userID)

		assert.NoError(t, err)
		assert.NotNil(t, cats)
		mockRepo.AssertExpectations(t)
	})

	t.Run("fail delete not found", func(t *testing.T) {
		catID := "cat-unknown"

		// Mock: Repo return error saat hapus
		mockRepo.On("Delete", ctx, catID, userID).Return(assert.AnError).Once()

		cats, err := svc.DeleteCategory(ctx, catID, userID)

		assert.Error(t, err)
		assert.Nil(t, cats)
	})
}