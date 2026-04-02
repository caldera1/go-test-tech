package usecase

import (
	"context"
	"task-api/internal/domain"
	"time"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) error
	FindByID(ctx context.Context, id string) (domain.User, error)
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	Update(ctx context.Context, user domain.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context) ([]domain.User, error)
}

type TaskRepository interface {
	Create(ctx context.Context, task domain.Task) error
	FindByID(ctx context.Context, id string) (domain.Task, error)
	Update(ctx context.Context, task domain.Task) error
	Delete(ctx context.Context, id string) error
	ListByAssignee(ctx context.Context, userID string) ([]domain.Task, error)
	ListAll(ctx context.Context) ([]domain.Task, error)
}

type CommentRepository interface {
	Create(ctx context.Context, comment domain.Comment) error
	ListByTask(ctx context.Context, taskID string) ([]domain.Comment, error)
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hash, password string) error
}

type TokenService interface {
	Generate(user domain.User) (TokenPair, error)
	Parse(ctx context.Context, token string) (Claims, error)
	ParseRefresh(ctx context.Context, token string) (Claims, error)
	Revoke(ctx context.Context, tokenID string) error
}

type Clock interface {
	Now() time.Time
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

type Claims struct {
	UserID             string
	Role               domain.Role
	TokenID            string
	MustChangePassword bool
}
