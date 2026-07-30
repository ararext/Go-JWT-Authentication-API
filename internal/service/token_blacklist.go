package service

import "context"

// TokenBlacklist is intentionally unimplemented for now — logout is client-side
// (the client just discards its tokens). This interface exists so a Redis-backed
// blacklist can be added later without changing AuthHandler or AuthService signatures.
type TokenBlacklist interface {
	Blacklist(ctx context.Context, token string) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
}