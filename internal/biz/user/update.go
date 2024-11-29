package user

import (
	"context"
	"github.com/coderlewin/kratosinit/api/proto/errcode"
	"github.com/coderlewin/kratosinit/internal/domain"
)

func (b *Biz) Update(ctx context.Context, ud *domain.User) error {
	_, err := b.FindById(ctx, ud.ID)
	if err != nil {
		return err
	}
	err = b.repo.Update(ctx, ud)
	if err != nil {
		return errcode.ErrorOperationFailed("更新失败")
	}
	return nil
}
