package usecase

import (
	"context"
	"errors"
	"task-api/internal/domain"
	"testing"
)

func TestDeleteUser_WithAssignedTasks_Blocked(t *testing.T) {
	userRepo := newInMemoryUserRepo()
	taskRepo := newInMemoryTaskRepo()
	uc := NewUserUseCase(userRepo, taskRepo, stubHasher{}, fixedClock{})

	userRepo.Create(context.Background(), domain.User{
		ID:    "exec-pedro",
		Email: "pedro.gonzalez@empresa.cl",
		Role:  domain.RoleExecutor,
	})
	taskRepo.Create(context.Background(), domain.Task{
		ID:             "tarea-pendiente",
		Title:          "Revisar logs de produccion",
		Status:         domain.StatusAssigned,
		AssignedUserID: "exec-pedro",
	})

	err := uc.Delete(context.Background(), "exec-pedro")
	if !errors.Is(err, domain.ErrUserHasAssignedTasks) {
		t.Fatalf("expected ErrUserHasAssignedTasks, got %v", err)
	}
}

func TestDeleteUser_WithoutTasks_OK(t *testing.T) {
	userRepo := newInMemoryUserRepo()
	taskRepo := newInMemoryTaskRepo()
	uc := NewUserUseCase(userRepo, taskRepo, stubHasher{}, fixedClock{})

	userRepo.Create(context.Background(), domain.User{
		ID:    "auditor-ana",
		Email: "ana.auditora@empresa.cl",
		Role:  domain.RoleAuditor,
	})

	err := uc.Delete(context.Background(), "auditor-ana")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	_, err = userRepo.FindByID(context.Background(), "auditor-ana")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("user should have been deleted, got %v", err)
	}
}
