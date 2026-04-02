package usecase

import (
	"context"
	"errors"
	"task-api/internal/domain"
	"testing"
	"time"
)

func setupTaskUseCase(now time.Time) (*TaskUseCase, *inMemoryTaskRepo, *inMemoryUserRepo) {
	taskRepo := newInMemoryTaskRepo()
	userRepo := newInMemoryUserRepo()
	commentRepo := newInMemoryCommentRepo()
	clock := fixedClock{t: now}
	uc := NewTaskUseCase(taskRepo, userRepo, commentRepo, clock)
	return uc, taskRepo, userRepo
}

func TestUpdateStatus_ExpiredTaskCannotTransition(t *testing.T) {
	now := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)
	uc, taskRepo, _ := setupTaskUseCase(now)

	taskRepo.Create(context.Background(), domain.Task{
		ID:             "tarea-vencida",
		Title:          "Revisar documentacion",
		Status:         domain.StatusAssigned,
		AssignedUserID: "exec-maria",
		DueDate:        now.Add(-24 * time.Hour), // vencio ayer
	})

	err := uc.UpdateStatus(context.Background(), "tarea-vencida", "exec-maria", domain.StatusStarted)
	if !errors.Is(err, domain.ErrForbidden) {
		t.Fatalf("expired task should be blocked by policy before reaching TransitionTo, got %v", err)
	}
}

func TestUpdateStatus_CannotSkipToFinalizado(t *testing.T) {
	now := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)
	uc, taskRepo, _ := setupTaskUseCase(now)

	taskRepo.Create(context.Background(), domain.Task{
		ID:             "tarea-nueva",
		Title:          "Implementar feature X",
		Status:         domain.StatusAssigned,
		AssignedUserID: "exec-maria",
		DueDate:        now.Add(24 * time.Hour),
	})

	// ASIGNADO -> FINALIZADO_EXITO directo no es valido, debe pasar por INICIADO
	err := uc.UpdateStatus(context.Background(), "tarea-nueva", "exec-maria", domain.StatusDoneOk)
	if !errors.Is(err, domain.ErrInvalidTaskTransition) {
		t.Fatalf("expected ErrInvalidTaskTransition, got %v", err)
	}
}

func TestUpdateStatus_ValidTransition(t *testing.T) {
	now := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)
	uc, taskRepo, _ := setupTaskUseCase(now)

	taskRepo.Create(context.Background(), domain.Task{
		ID:             "tarea-activa",
		Title:          "Migrar base de datos",
		Status:         domain.StatusAssigned,
		AssignedUserID: "exec-maria",
		DueDate:        now.Add(7 * 24 * time.Hour), // vence en una semana
	})

	err := uc.UpdateStatus(context.Background(), "tarea-activa", "exec-maria", domain.StatusStarted)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updated, err := taskRepo.FindByID(context.Background(), "tarea-activa")
	if err != nil {
		t.Fatal(err)
	}
	if updated.Status != domain.StatusStarted {
		t.Fatalf("expected status INICIADO, got %s", updated.Status)
	}
}

func TestCreateTask_RejectsAssigneeWithoutExecutorRole(t *testing.T) {
	now := time.Date(2025, 6, 1, 12, 0, 0, 0, time.UTC)
	uc, _, userRepo := setupTaskUseCase(now)

	userRepo.Create(context.Background(), domain.User{
		ID:    "auditor-carlos",
		Email: "carlos.audit@empresa.cl",
		Role:  domain.RoleAuditor,
	})

	_, err := uc.Create(context.Background(), "Deploy a produccion", "Desplegar v2.0", now.Add(24*time.Hour), "auditor-carlos", "admin-01")
	if !errors.Is(err, domain.ErrAssigneeNotExecutor) {
		t.Fatalf("expected ErrAssigneeNotExecutor, got %v", err)
	}
}
