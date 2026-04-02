package usecase

import (
	"context"
	"task-api/internal/domain"
)

type LoginResult struct {
	Tokens             TokenPair
	MustChangePassword bool
}

type AuthUseCase struct {
	users  UserRepository
	hasher PasswordHasher
	tokens TokenService
}

func NewAuthUseCase(users UserRepository, hasher PasswordHasher, tokens TokenService) *AuthUseCase {
	return &AuthUseCase{users: users, hasher: hasher, tokens: tokens}
}

func (uc *AuthUseCase) Login(ctx context.Context, email, password string) (LoginResult, error) {
	user, err := uc.users.FindByEmail(ctx, email)
	if err != nil {
		return LoginResult{}, domain.ErrInvalidCredentials
	}

	if err := uc.hasher.Compare(user.PasswordHash, password); err != nil {
		return LoginResult{}, domain.ErrInvalidCredentials
	}

	tokens, err := uc.tokens.Generate(user)
	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{
		Tokens:             tokens,
		MustChangePassword: user.MustChangePassword,
	}, nil
}

func (uc *AuthUseCase) ChangePassword(ctx context.Context, userID, currentPass, newPass string) error {
	user, err := uc.users.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := uc.hasher.Compare(user.PasswordHash, currentPass); err != nil {
		return domain.ErrInvalidCredentials
	}

	hash, err := uc.hasher.Hash(newPass)
	if err != nil {
		return err
	}

	user.PasswordHash = hash
	user.MustChangePassword = false
	return uc.users.Update(ctx, user)
}

func (uc *AuthUseCase) Logout(ctx context.Context, refreshToken string) error {
	// Revocamos el refresh token. El access token expira solo (TTL 15min).
	// Trade-off: ventana máxima de 15min post-logout sin blacklist de access tokens.
	claims, err := uc.tokens.ParseRefresh(ctx, refreshToken)
	if err != nil {
		return err
	}
	return uc.tokens.Revoke(ctx, claims.TokenID)
}
