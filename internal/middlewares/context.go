package middlewares

import (
	"context"
	"github.com/golang-jwt/jwt/v4"
)

type contextKey string

const userCtxKey = contextKey("userClaims")

func WithClaims(ctx context.Context, claims jwt.MapClaims) context.Context {
	return context.WithValue(ctx, userCtxKey, claims)
}

func FromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(userCtxKey).(jwt.MapClaims)
	if !ok {
		return &Claims{}, false
	}

	chatIDFloat, ok := claims["id"].(float64)
	if !ok {
		return &Claims{}, false
	}
	chatID := int64(chatIDFloat)

	return &Claims{
		ChatID: chatID,
	}, ok
}

type Claims struct {
	ChatID int64
}
