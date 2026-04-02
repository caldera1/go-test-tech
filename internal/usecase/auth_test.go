package usecase

import (
	"context"
	"errors"
	"strings"
	"task-api/internal/domain"
	"testing"
)

type stubHasher struct{}

func (s stubHasher) Hash(p string) (string, error) { return "hashed_" + p, nil }
func (s stubHasher) Compare(hash, p string) error {
	if hash != "hashed_"+p {
		return domain.ErrInvalidCredentials
	}
	return nil
}

type stubTokenService struct {
	generated int
	revoked   []string
}

func (s *stubTokenService) Generate(u domain.User) (TokenPair, error) {
	s.generated++
	return TokenPair{AccessToken: "access_" + u.ID, RefreshToken: "refresh_" + u.ID}, nil
}

func (s *stubTokenService) Parse(ctx context.Context, token string) (Claims, error) {
	return Claims{}, domain.ErrInvalidToken
}

func (s *stubTokenService) ParseRefresh(ctx context.Context, token string) (Claims, error) {
	if !strings.HasPrefix(token, "refresh_") {
		return Claims{}, domain.ErrInvalidToken
	}
	userID := strings.TrimPrefix(token, "refresh_")
	return Claims{TokenID: "jti_" + userID}, nil
}

func (s *stubTokenService) Revoke(ctx context.Context, tokenID string) error {
	s.revoked = append(s.revoked, tokenID)
	return nil
}

func TestLogin_OK(t *testing.T) {
	repo := newInMemoryUserRepo()
	repo.Create(context.Background(), domain.User{
		ID:                 "exec-01",
		Email:              "maria.lopez@empresa.cl",
		PasswordHash:       "hashed_temporal123",
		Role:               domain.RoleExecutor,
		MustChangePassword: true,
	})

	tokens := &stubTokenService{}
	uc := NewAuthUseCase(repo, stubHasher{}, tokens)

	result, err := uc.Login(context.Background(), "maria.lopez@empresa.cl", "temporal123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !result.MustChangePassword {
		t.Fatal("first login should flag MustChangePassword=true")
	}
	if tokens.generated != 1 {
		t.Fatalf("expected 1 token pair generated, got %d", tokens.generated)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	repo := newInMemoryUserRepo()
	repo.Create(context.Background(), domain.User{
		ID:           "exec-01",
		Email:        "maria.lopez@empresa.cl",
		PasswordHash: "hashed_temporal123",
	})

	uc := NewAuthUseCase(repo, stubHasher{}, &stubTokenService{})

	_, err := uc.Login(context.Background(), "maria.lopez@empresa.cl", "incorrecto")
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogin_UserNotFound_ReturnsGenericError(t *testing.T) {
	repo := newInMemoryUserRepo()
	uc := NewAuthUseCase(repo, stubHasher{}, &stubTokenService{})

	// No debe revelar si el email existe o no — siempre ErrInvalidCredentials
	_, err := uc.Login(context.Background(), "inexistente@empresa.cl", "cualquiera")
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials (not ErrNotFound), got %v", err)
	}
}

func TestChangePassword_OK(t *testing.T) {
	repo := newInMemoryUserRepo()
	repo.Create(context.Background(), domain.User{
		ID:                 "exec-01",
		Email:              "maria.lopez@empresa.cl",
		PasswordHash:       "hashed_temporal123",
		MustChangePassword: true,
	})

	uc := NewAuthUseCase(repo, stubHasher{}, &stubTokenService{})

	err := uc.ChangePassword(context.Background(), "exec-01", "temporal123", "NuevaSegura456!")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	updated, _ := repo.FindByID(context.Background(), "exec-01")
	if updated.MustChangePassword {
		t.Fatal("MustChangePassword should be false after password change")
	}
	if updated.PasswordHash != "hashed_NuevaSegura456!" {
		t.Fatalf("expected updated hash, got '%s'", updated.PasswordHash)
	}
}

func TestChangePassword_WrongCurrentPassword(t *testing.T) {
	repo := newInMemoryUserRepo()
	repo.Create(context.Background(), domain.User{
		ID:           "exec-01",
		Email:        "maria.lopez@empresa.cl",
		PasswordHash: "hashed_temporal123",
	})

	uc := NewAuthUseCase(repo, stubHasher{}, &stubTokenService{})

	err := uc.ChangePassword(context.Background(), "exec-01", "equivocada", "NuevaSegura456!")
	if !errors.Is(err, domain.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestLogout_OK(t *testing.T) {
	tokens := &stubTokenService{}
	uc := NewAuthUseCase(newInMemoryUserRepo(), stubHasher{}, tokens)

	err := uc.Logout(context.Background(), "refresh_user1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(tokens.revoked) != 1 || tokens.revoked[0] != "jti_user1" {
		t.Fatalf("expected revoked=[jti_user1], got %v", tokens.revoked)
	}
}
