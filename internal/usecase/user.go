package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"task-api/internal/domain"

	"github.com/google/uuid"
)

type CreateUserResult struct {
	User              domain.User
	TemporaryPassword string
}

type UserUseCase struct {
	users  UserRepository
	tasks  TaskRepository
	hasher PasswordHasher
	clock  Clock
}

func NewUserUseCase(users UserRepository, tasks TaskRepository, hasher PasswordHasher, clock Clock) *UserUseCase {
	return &UserUseCase{users: users, tasks: tasks, hasher: hasher, clock: clock}
}

func (uc *UserUseCase) Create(ctx context.Context, email string, role domain.Role) (CreateUserResult, error) {
	if role == domain.RoleAdmin {
		return CreateUserResult{}, domain.ErrAdminCannotBeCreated
	}

	// La contraseña temporal se devuelve solo en este momento.
	// Después del hash no es recuperable — en producción se enviaría por canal seguro.
	tempPass, err := generateTemporaryPassword()
	if err != nil {
		return CreateUserResult{}, err
	}

	hash, err := uc.hasher.Hash(tempPass)
	if err != nil {
		return CreateUserResult{}, err
	}

	user := domain.User{
		ID:                 uuid.New().String(),
		Email:              email,
		PasswordHash:       hash,
		Role:               role,
		MustChangePassword: true,
		CreatedAt:          uc.clock.Now(),
	}

	if err := uc.users.Create(ctx, user); err != nil {
		return CreateUserResult{}, err
	}

	return CreateUserResult{User: user, TemporaryPassword: tempPass}, nil
}

func (uc *UserUseCase) GetByID(ctx context.Context, id string) (domain.User, error) {
	return uc.users.FindByID(ctx, id)
}

func (uc *UserUseCase) Update(ctx context.Context, id, email string, role domain.Role) (domain.User, error) {
	user, err := uc.users.FindByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}

	if role == domain.RoleAdmin {
		return domain.User{}, domain.ErrAdminCannotBeCreated
	}

	user.Email = email
	user.Role = role

	if err := uc.users.Update(ctx, user); err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (uc *UserUseCase) Delete(ctx context.Context, id string) error {
	tasks, err := uc.tasks.ListByAssignee(ctx, id)
	if err != nil {
		return err
	}
	if len(tasks) > 0 {
		return domain.ErrUserHasAssignedTasks
	}
	return uc.users.Delete(ctx, id)
}

func (uc *UserUseCase) List(ctx context.Context) ([]domain.User, error) {
	return uc.users.List(ctx)
}

func generateTemporaryPassword() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
