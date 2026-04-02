package jwt

import "context"

type RevokedTokenRepository interface {
	Revoke(ctx context.Context, tokenID string) error
	IsRevoked(ctx context.Context, tokenID string) (bool, error)
}
