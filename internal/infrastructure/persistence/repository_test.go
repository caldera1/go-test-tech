package persistence

import (
	"context"
	"errors"
	"task-api/internal/domain"
	"task-api/internal/infrastructure/persistence/models"
	"testing"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	if err := db.AutoMigrate(
		&models.User{},
		&models.Task{},
		&models.Comment{},
		&models.RevokedToken{},
	); err != nil {
		t.Fatalf("failed to migrate: %v", err)
	}

	return db
}

func TestUserRepository_CreateAndFind(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)
	ctx := context.Background()

	user := domain.User{
		ID:                 "exec-maria-01",
		Email:              "maria.lopez@empresa.cl",
		PasswordHash:       "$2a$10$fakehash",
		Role:               domain.RoleExecutor,
		MustChangePassword: true,
		CreatedAt:          time.Now(),
	}

	if err := repo.Create(ctx, user); err != nil {
		t.Fatalf("create failed: %v", err)
	}

	found, err := repo.FindByID(ctx, "exec-maria-01")
	if err != nil {
		t.Fatalf("find by id failed: %v", err)
	}
	if found.Email != "maria.lopez@empresa.cl" {
		t.Fatalf("expected maria.lopez@empresa.cl, got %s", found.Email)
	}

	foundByEmail, err := repo.FindByEmail(ctx, "maria.lopez@empresa.cl")
	if err != nil {
		t.Fatalf("find by email failed: %v", err)
	}
	if foundByEmail.ID != "exec-maria-01" {
		t.Fatalf("expected exec-maria-01, got %s", foundByEmail.ID)
	}
}

func TestUserRepository_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewUserRepository(db)

	_, err := repo.FindByID(context.Background(), "nonexistent")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestTaskRepository_CreateAndFind(t *testing.T) {
	db := setupTestDB(t)
	userRepo := NewUserRepository(db)
	taskRepo := NewTaskRepository(db)
	ctx := context.Background()

	executor := domain.User{
		ID:           "exec-01",
		Email:        "pedro.ejecutor@empresa.cl",
		PasswordHash: "$2a$10$fakehash",
		Role:         domain.RoleExecutor,
	}
	userRepo.Create(ctx, executor)

	admin := domain.User{
		ID:           "admin-01",
		Email:        "admin@empresa.cl",
		PasswordHash: "$2a$10$fakehash",
		Role:         domain.RoleAdmin,
	}
	userRepo.Create(ctx, admin)

	task := domain.Task{
		ID:              "tarea-deploy",
		Title:           "Deploy v2.1 a staging",
		Description:     "Ejecutar pipeline de deploy para ambiente de staging",
		DueDate:         time.Now().Add(24 * time.Hour),
		Status:          domain.StatusAssigned,
		AssignedUserID:  "exec-01",
		CreatedByUserID: "admin-01",
	}

	if err := taskRepo.Create(ctx, task); err != nil {
		t.Fatalf("create task failed: %v", err)
	}

	found, err := taskRepo.FindByID(ctx, "tarea-deploy")
	if err != nil {
		t.Fatalf("find task failed: %v", err)
	}
	if found.Title != "Deploy v2.1 a staging" {
		t.Fatalf("expected 'Deploy v2.1 a staging', got %s", found.Title)
	}
}

func TestTaskRepository_DeleteNotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTaskRepository(db)

	err := repo.Delete(context.Background(), "nonexistent")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
