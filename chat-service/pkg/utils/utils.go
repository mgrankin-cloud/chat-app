package utils

import (
	"context"
	"errors"
)

type contextKey string

const userIDKey contextKey = "userID"

var ErrUserIDNotFound = errors.New("user ID not found in context")

func WithUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userIDKey, userID)
}

func GetUserIDFromContext(ctx context.Context) (int64, error) {
	userID, ok := ctx.Value(userIDKey).(int64)
	if !ok {
		return 0, ErrUserIDNotFound
	}
	return userID, nil
}