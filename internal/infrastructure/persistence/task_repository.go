package persistence

import (
	"context"
	"errors"
	"task-api/internal/domain"
	"task-api/internal/infrastructure/persistence/models"

	"gorm.io/gorm"
)

type TaskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(ctx context.Context, task domain.Task) error {
	m := models.TaskToModel(task)
	return r.db.WithContext(ctx).Create(&m).Error
}

func (r *TaskRepository) FindByID(ctx context.Context, id string) (domain.Task, error) {
	var m models.Task
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return domain.Task{}, domain.ErrNotFound
		}
		return domain.Task{}, err
	}
	return models.TaskToDomain(m), nil
}

func (r *TaskRepository) Update(ctx context.Context, task domain.Task) error {
	m := models.TaskToModel(task)
	result := r.db.WithContext(ctx).Model(&models.Task{}).Where("id = ?", m.ID).Select("*").Updates(m)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *TaskRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).Delete(&models.Task{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return domain.ErrNotFound
	}
	return nil
}

func (r *TaskRepository) ListByAssignee(ctx context.Context, userID string) ([]domain.Task, error) {
	var ms []models.Task
	if err := r.db.WithContext(ctx).Where("assigned_user_id = ?", userID).Find(&ms).Error; err != nil {
		return nil, err
	}
	tasks := make([]domain.Task, len(ms))
	for i, m := range ms {
		tasks[i] = models.TaskToDomain(m)
	}
	return tasks, nil
}

func (r *TaskRepository) ListAll(ctx context.Context) ([]domain.Task, error) {
	var ms []models.Task
	if err := r.db.WithContext(ctx).Find(&ms).Error; err != nil {
		return nil, err
	}
	tasks := make([]domain.Task, len(ms))
	for i, m := range ms {
		tasks[i] = models.TaskToDomain(m)
	}
	return tasks, nil
}
