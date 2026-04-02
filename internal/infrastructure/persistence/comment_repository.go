package persistence

import (
	"context"
	"task-api/internal/domain"
	"task-api/internal/infrastructure/persistence/models"

	"gorm.io/gorm"
)

type CommentRepository struct {
	db *gorm.DB
}

func NewCommentRepository(db *gorm.DB) *CommentRepository {
	return &CommentRepository{db: db}
}

func (r *CommentRepository) Create(ctx context.Context, comment domain.Comment) error {
	m := models.CommentToModel(comment)
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *CommentRepository) ListByTask(ctx context.Context, taskID string) ([]domain.Comment, error) {
	var ms []models.Comment
	if err := r.db.WithContext(ctx).Where("task_id = ?", taskID).Order("created_at ASC").Find(&ms).Error; err != nil {
		return nil, err
	}
	comments := make([]domain.Comment, len(ms))
	for i, m := range ms {
		comments[i] = models.CommentToDomain(m)
	}
	return comments, nil
}
