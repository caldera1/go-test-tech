package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"task-api/internal/domain"
	httphandler "task-api/internal/handler/http"
	"task-api/internal/infrastructure/crypto"
	jwtinfra "task-api/internal/infrastructure/jwt"
	"task-api/internal/infrastructure/persistence"
	"task-api/internal/usecase"
	"time"
)

type config struct {
	dsn       string
	jwtSecret string
	port      string
}

func loadConfig() config {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return config{dsn: dsn, jwtSecret: secret, port: port}
}

func main() {
	cfg := loadConfig()

	db, err := persistence.Connect(cfg.dsn)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	userRepo := persistence.NewUserRepository(db)
	taskRepo := persistence.NewTaskRepository(db)
	commentRepo := persistence.NewCommentRepository(db)
	revokedTokenRepo := persistence.NewRevokedTokenRepository(db)

	hasher := crypto.BcryptHasher{}
	clock := usecase.RealClock{}
	jwtSvc := jwtinfra.NewJWTService(cfg.jwtSecret, revokedTokenRepo)

	seedAdmin(userRepo, hasher)

	authUC := usecase.NewAuthUseCase(userRepo, hasher, jwtSvc)
	userUC := usecase.NewUserUseCase(userRepo, taskRepo, hasher, clock)
	taskUC := usecase.NewTaskUseCase(taskRepo, userRepo, commentRepo, clock)

	router := httphandler.SetupRouter(authUC, userUC, taskUC, jwtSvc)

	srv := &http.Server{
		Addr:         ":" + cfg.port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	log.Printf("server starting on port %s", cfg.port)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("server error: %v", err)
	}
}

func seedAdmin(users *persistence.UserRepository, hasher crypto.BcryptHasher) {
	// Crea el admin inicial solo si no existe. Solo para bootstrap.
	// En producción las credenciales iniciales deberían venir por secret management.
	ctx := context.Background()
	_, err := users.FindByEmail(ctx, "admin@admin.com")
	if err == nil {
		return
	}
	if !errors.Is(err, domain.ErrNotFound) {
		log.Fatalf("failed to check admin user: %v", err)
	}

	hash, err := hasher.Hash("Admin1234!")
	if err != nil {
		log.Fatalf("failed to hash admin password: %v", err)
	}

	admin := domain.User{
		ID:                 "00000000-0000-0000-0000-000000000001",
		Email:              "admin@admin.com",
		PasswordHash:       hash,
		Role:               domain.RoleAdmin,
		MustChangePassword: true,
	}

	if err := users.Create(ctx, admin); err != nil {
		log.Fatalf("failed to seed admin user: %v", err)
	}

	log.Println("admin user seeded: admin@admin.com / Admin1234!")
}
