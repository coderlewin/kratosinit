package ctxutils

import "context"

func MustGetUserId(ctx context.Context) int64 {
	return int64(ctx.Value("userId").(float64))
}
