package persistence

import (
	"context"
	"task-api/internal/infrastructure/persistence/models"
	"time"

	"gorm.io/gorm"
)

type RevokedTokenRepository struct {
	db *gorm.DB
}

func NewRevokedTokenRepository(db *gorm.DB) *RevokedTokenRepository {
	return &RevokedTokenRepository{db: db}
}

func (r *RevokedTokenRepository) Revoke(ctx context.Context, tokenID string) error {
	token := models.RevokedToken{ID: tokenID, RevokedAt: time.Now()}
	return r.db.WithContext(ctx).FirstOrCreate(&token, "id = ?", tokenID).Error
}

func (r *RevokedTokenRepository) IsRevoked(ctx context.Context, tokenID string) (bool, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&models.RevokedToken{}).Where("id = ?", tokenID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
