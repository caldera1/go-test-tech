package persistence

import (
	"context"
	"errors"
	"task-api/internal/domain"
	"task-api/internal/infrastructure/persistence/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) error {
	m := models.UserToModel(user)
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (domain.User, error) {
	var m models.User
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, err
	}
	return models.UserToDomain(m), nil
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	var m models.User
	if err := r.db.WithContext(ctx).First(&m, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.User{}, domain.ErrNotFound
		}
		return domain.User{}, err
	}
	return models.UserToDomain(m), nil
}

func (r *UserRepository) Update(ctx context.Context, user domain.User) error {
	m := models.UserToModel(user)
	result := r.db.WithContext(ctx).Model(&models.User{}).Where("id = ?", m.ID).Select("*").Updates(m)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *UserRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.User{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *UserRepository) List(ctx context.Context) ([]domain.User, error) {
	var ms []models.User
	if err := r.db.WithContext(ctx).Find(&ms).Error; err != nil {
		return nil, err
	}
	users := make([]domain.User, len(ms))
	for i, m := range ms {
		users[i] = models.UserToDomain(m)
	}
	return users, nil
}
