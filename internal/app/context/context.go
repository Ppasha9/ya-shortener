package context

import (
	"context"
	"errors"
)

type ctxKey int

const (
	userIDKey ctxKey = iota
)

func WithUserID(ctx context.Context, userID uint32) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetUserID(ctx context.Context) (uint32, error) {
	userID := ctx.Value(userIDKey)
	if userID == nil {
		return 0, errors.New("failed to get user id from context")
	}

	return userID.(uint32), nil
}
