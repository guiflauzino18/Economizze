package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/guiflauzino18/economizze/internal/domain/entities"
	errorsDomain "github.com/guiflauzino18/economizze/internal/domain/errors"
	"github.com/guiflauzino18/economizze/internal/ports"
	"gorm.io/gorm"
)

var _ ports.CategoryRepository = (*categoryRepository)(nil)

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) ports.CategoryRepository {
	return categoryRepository{db}
}

// FindAvailableForUser implements [ports.CategoryRepository].
func (c categoryRepository) FindAvailableForUser(ctx context.Context, userID uuid.UUID) ([]*entities.Category, error) {
	var models []CategoryModel

	err := c.db.WithContext(ctx).
		Where("(user_id IS NULL OR user_id = ?) AND active = true", userID).
		Order("user_id IS NOT NULL, name ASC").
		Find(&models).Error

	if err != nil {
		return nil, fmt.Errorf("categoryRepository.FindAvailableForUser: %w", err)
	}

	categories := make([]*entities.Category, 0, len(models))

	for _, m := range models {
		category := modelToCategory(m)
		categories = append(categories, category)
	}

	return categories, nil

}

// FindByID implements [ports.CategoryRepository].
func (c categoryRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Category, error) {
	var model CategoryModel

	err := c.db.WithContext(ctx).
		Where("id = ? AND active = true", id).
		First(&model).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("categoryRepository.FindByID %s: %w", id, errorsDomain.ErrNotFound)
	}

	if err != nil {
		return nil, fmt.Errorf("categoryRepository.FindByID: %w", err)
	}

	return modelToCategory(model), nil
}

// Save implements [ports.CategoryRepository].
func (c categoryRepository) Save(ctx context.Context, cat *entities.Category) error {

	model := CategoryToModel(cat)

	err := c.db.WithContext(ctx).
		Save(&model).Error

	if err != nil {
		return fmt.Errorf("categoryRepository.Save: %w", err)
	}

	return nil
}
