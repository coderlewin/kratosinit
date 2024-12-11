package ctxutils

import "context"

type (
	userKey        struct{}
	accessTokenKey struct{}
)

func NewUserID(ctx context.Context, userID int64) context.Context {
	return context.WithValue(ctx, userKey{}, userID)
}

func FromUserID(ctx context.Context) int64 {
	userID, _ := ctx.Value(userKey{}).(int64)
	return userID
}

// NewAccessToken put accessToken into context.
func NewAccessToken(ctx context.Context, accessToken string) context.Context {
	return context.WithValue(ctx, accessTokenKey{}, accessToken)
}

// FromAccessToken extract accessToken from context.
func FromAccessToken(ctx context.Context) string {
	accessToken, _ := ctx.Value(accessTokenKey{}).(string)
	return accessToken
}
