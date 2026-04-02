package usecase

import (
	"context"
	"sync"
	"task-api/internal/domain"
	"time"
)

type fixedClock struct {
	t time.Time
}

func (c fixedClock) Now() time.Time { return c.t }

type inMemoryUserRepo struct {
	mu    sync.RWMutex
	users map[string]domain.User
}

func newInMemoryUserRepo() *inMemoryUserRepo {
	return &inMemoryUserRepo{users: make(map[string]domain.User)}
}

func (r *inMemoryUserRepo) Create(_ context.Context, user domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.users[user.ID] = user
	return nil
}

func (r *inMemoryUserRepo) FindByID(_ context.Context, id string) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok {
		return domain.User{}, domain.ErrNotFound
	}
	return u, nil
}

func (r *inMemoryUserRepo) FindByEmail(_ context.Context, email string) (domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.users {
		if u.Email == email {
			return u, nil
		}
	}
	return domain.User{}, domain.ErrNotFound
}

func (r *inMemoryUserRepo) Update(_ context.Context, user domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[user.ID]; !ok {
		return domain.ErrNotFound
	}
	r.users[user.ID] = user
	return nil
}

func (r *inMemoryUserRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.users[id]; !ok {
		return domain.ErrNotFound
	}
	delete(r.users, id)
	return nil
}

func (r *inMemoryUserRepo) List(_ context.Context) ([]domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]domain.User, 0, len(r.users))
	for _, u := range r.users {
		result = append(result, u)
	}
	return result, nil
}

type inMemoryTaskRepo struct {
	mu    sync.RWMutex
	tasks map[string]domain.Task
}

func newInMemoryTaskRepo() *inMemoryTaskRepo {
	return &inMemoryTaskRepo{tasks: make(map[string]domain.Task)}
}

func (r *inMemoryTaskRepo) Create(_ context.Context, task domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.tasks[task.ID] = task
	return nil
}

func (r *inMemoryTaskRepo) FindByID(_ context.Context, id string) (domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	t, ok := r.tasks[id]
	if !ok {
		return domain.Task{}, domain.ErrNotFound
	}
	return t, nil
}

func (r *inMemoryTaskRepo) Update(_ context.Context, task domain.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tasks[task.ID]; !ok {
		return domain.ErrNotFound
	}
	r.tasks[task.ID] = task
	return nil
}

func (r *inMemoryTaskRepo) Delete(_ context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.tasks[id]; !ok {
		return domain.ErrNotFound
	}
	delete(r.tasks, id)
	return nil
}

func (r *inMemoryTaskRepo) ListByAssignee(_ context.Context, userID string) ([]domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Task
	for _, t := range r.tasks {
		if t.AssignedUserID == userID {
			result = append(result, t)
		}
	}
	return result, nil
}

func (r *inMemoryTaskRepo) ListAll(_ context.Context) ([]domain.Task, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]domain.Task, 0, len(r.tasks))
	for _, t := range r.tasks {
		result = append(result, t)
	}
	return result, nil
}

type inMemoryCommentRepo struct {
	mu       sync.RWMutex
	comments []domain.Comment
}

func newInMemoryCommentRepo() *inMemoryCommentRepo {
	return &inMemoryCommentRepo{}
}

func (r *inMemoryCommentRepo) Create(_ context.Context, comment domain.Comment) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.comments = append(r.comments, comment)
	return nil
}

func (r *inMemoryCommentRepo) ListByTask(_ context.Context, taskID string) ([]domain.Comment, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []domain.Comment
	for _, c := range r.comments {
		if c.TaskID == taskID {
			result = append(result, c)
		}
	}
	return result, nil
}
