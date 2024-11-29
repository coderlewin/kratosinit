package user

import (
	"context"
	"github.com/coderlewin/kratosinit/api/proto/errcode"
)

func (b *Biz) Delete(ctx context.Context, id int64) error {
	_, err := b.FindById(ctx, id)
	if err != nil {
		return err
	}
	err = b.repo.Delete(ctx, id)
	if err != nil {
		b.log.Errorf("delete user error: %v", err)
		return errcode.ErrorOperationFailed("删除失败")
	}
	return nil
}
