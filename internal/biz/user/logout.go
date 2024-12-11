package user

import (
	"context"
	"github.com/coderlewin/kratosinit/api/proto/errcode"
)

func (b *Biz) Logout(ctx context.Context, accessToken string) error {
	err := b.authn.Destroy(ctx, accessToken)
	if err != nil {
		return errcode.ErrorUnknown("Failed to logout")
	}
	return nil
}
