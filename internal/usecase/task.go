package usecase

import (
	"context"
	"task-api/internal/domain"
	"time"

	"github.com/google/uuid"
)

type TaskUseCase struct {
	tasks    TaskRepository
	users    UserRepository
	comments CommentRepository
	clock    Clock
}

func NewTaskUseCase(tasks TaskRepository, users UserRepository, comments CommentRepository, clock Clock) *TaskUseCase {
	return &TaskUseCase{tasks: tasks, users: users, comments: comments, clock: clock}
}

func (uc *TaskUseCase) Create(ctx context.Context, title, description string, dueDate time.Time, assignedUserID, createdByUserID string) (domain.Task, error) {
	assignee, err := uc.users.FindByID(ctx, assignedUserID)
	if err != nil {
		return domain.Task{}, err
	}

	if assignee.Role != domain.RoleExecutor {
		return domain.Task{}, domain.ErrAssigneeNotExecutor
	}

	task := domain.Task{
		ID:              uuid.New().String(),
		Title:           title,
		Description:     description,
		DueDate:         dueDate,
		Status:          domain.StatusAssigned,
		AssignedUserID:  assignedUserID,
		CreatedByUserID: createdByUserID,
		CreatedAt:       uc.clock.Now(),
	}

	if err := uc.tasks.Create(ctx, task); err != nil {
		return domain.Task{}, err
	}

	return task, nil
}

func (uc *TaskUseCase) AdminUpdate(ctx context.Context, taskID, title, description string, dueDate time.Time) error {
	task, err := uc.tasks.FindByID(ctx, taskID)
	if err != nil {
		return err
	}

	if !task.CanBeAdminModified() {
		return domain.ErrForbidden
	}

	task.Title = title
	task.Description = description
	task.DueDate = dueDate

	return uc.tasks.Update(ctx, task)
}

func (uc *TaskUseCase) AdminDelete(ctx context.Context, taskID string) error {
	task, err := uc.tasks.FindByID(ctx, taskID)
	if err != nil {
		return err
	}

	if !task.CanBeAdminModified() {
		return domain.ErrForbidden
	}

	return uc.tasks.Delete(ctx, taskID)
}

func (uc *TaskUseCase) UpdateStatus(ctx context.Context, taskID, userID string, next domain.TaskStatus) error {
	task, err := uc.tasks.FindByID(ctx, taskID)
	if err != nil {
		return err
	}

	now := uc.clock.Now()
	if !domain.CanExecutorUpdateTask(task, userID, now) {
		return domain.ErrForbidden
	}

	if err := task.TransitionTo(next, now); err != nil {
		return err
	}

	return uc.tasks.Update(ctx, task)
}

func (uc *TaskUseCase) AddComment(ctx context.Context, taskID, userID, body string) error {
	task, err := uc.tasks.FindByID(ctx, taskID)
	if err != nil {
		return err
	}

	if !domain.CanAddComment(task, userID, uc.clock.Now()) {
		return domain.ErrForbidden
	}

	comment := domain.Comment{
		ID:        uuid.New().String(),
		TaskID:    taskID,
		AuthorID:  userID,
		Body:      body,
		CreatedAt: uc.clock.Now(),
	}

	return uc.comments.Create(ctx, comment)
}

type TaskDetail struct {
	Task     domain.Task
	Comments []domain.Comment
}

func (uc *TaskUseCase) GetDetail(ctx context.Context, taskID, userID string, role domain.Role) (TaskDetail, error) {
	task, err := uc.tasks.FindByID(ctx, taskID)
	if err != nil {
		return TaskDetail{}, err
	}

	if !domain.CanViewTask(task, userID, role) {
		return TaskDetail{}, domain.ErrForbidden
	}

	comments, err := uc.comments.ListByTask(ctx, taskID)
	if err != nil {
		return TaskDetail{}, err
	}

	return TaskDetail{Task: task, Comments: comments}, nil
}

func (uc *TaskUseCase) GetMine(ctx context.Context, taskID, userID string) (TaskDetail, error) {
	task, err := uc.tasks.FindByID(ctx, taskID)
	if err != nil {
		return TaskDetail{}, err
	}

	if task.AssignedUserID != userID {
		return TaskDetail{}, domain.ErrForbidden
	}

	comments, err := uc.comments.ListByTask(ctx, taskID)
	if err != nil {
		return TaskDetail{}, err
	}

	return TaskDetail{Task: task, Comments: comments}, nil
}

func (uc *TaskUseCase) ListMine(ctx context.Context, userID string) ([]domain.Task, error) {
	return uc.tasks.ListByAssignee(ctx, userID)
}

func (uc *TaskUseCase) ListAll(ctx context.Context) ([]domain.Task, error) {
	return uc.tasks.ListAll(ctx)
}
