package jwt

import (
	"context"
	"task-api/internal/domain"
	"task-api/internal/usecase"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"

	accessTokenTTL  = 15 * time.Minute
	refreshTokenTTL = 7 * 24 * time.Hour
)

type Claims struct {
	UserID             string    `json:"user_id"`
	Role               string    `json:"role"`
	TokenID            string    `json:"token_id"`
	TokenType          TokenType `json:"typ"`
	MustChangePassword bool      `json:"must_change_password"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secret      []byte
	revokedRepo RevokedTokenRepository
}

func NewJWTService(secret string, revokedRepo RevokedTokenRepository) *JWTService {
	return &JWTService{
		secret:      []byte(secret),
		revokedRepo: revokedRepo,
	}
}

func (s *JWTService) Generate(user domain.User) (usecase.TokenPair, error) {
	now := time.Now()

	accessID := uuid.New().String()
	accessClaims := Claims{
		UserID:             user.ID,
		Role:               string(user.Role),
		TokenID:            accessID,
		TokenType:          TokenTypeAccess,
		MustChangePassword: user.MustChangePassword,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(accessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        accessID,
		},
	}

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(s.secret)
	if err != nil {
		return usecase.TokenPair{}, err
	}

	refreshID := uuid.New().String()
	refreshClaims := Claims{
		UserID:             user.ID,
		Role:               string(user.Role),
		TokenID:            refreshID,
		TokenType:          TokenTypeRefresh,
		MustChangePassword: user.MustChangePassword,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(refreshTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        refreshID,
		},
	}

	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(s.secret)
	if err != nil {
		return usecase.TokenPair{}, err
	}

	return usecase.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *JWTService) Parse(ctx context.Context, tokenStr string) (usecase.Claims, error) {
	return s.parseWithType(ctx, tokenStr, TokenTypeAccess)
}

func (s *JWTService) ParseRefresh(ctx context.Context, tokenStr string) (usecase.Claims, error) {
	return s.parseWithType(ctx, tokenStr, TokenTypeRefresh)
}

func (s *JWTService) parseWithType(ctx context.Context, tokenStr string, expectedType TokenType) (usecase.Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrInvalidToken
		}
		return s.secret, nil
	})
	if err != nil {
		return usecase.Claims{}, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return usecase.Claims{}, domain.ErrInvalidToken
	}

	// Validamos el tipo para evitar que un refresh token sea usado como access token
	// y viceversa. Ambos tienen la misma firma, solo el claim typ los diferencia.
	if claims.TokenType != expectedType {
		return usecase.Claims{}, domain.ErrInvalidToken
	}

	revoked, err := s.revokedRepo.IsRevoked(ctx, claims.TokenID)
	if err != nil {
		return usecase.Claims{}, err
	}
	if revoked {
		return usecase.Claims{}, domain.ErrInvalidToken
	}

	return usecase.Claims{
		UserID:             claims.UserID,
		Role:               domain.Role(claims.Role),
		TokenID:            claims.TokenID,
		MustChangePassword: claims.MustChangePassword,
	}, nil
}

func (s *JWTService) Revoke(ctx context.Context, tokenID string) error {
	return s.revokedRepo.Revoke(ctx, tokenID)
}
