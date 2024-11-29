package user

import (
	"context"
	"github.com/coderlewin/kratosinit/api/proto/errcode"
	"github.com/coderlewin/kratosinit/internal/domain"
)

func (b *Biz) FindById(ctx context.Context, id int64) (*domain.User, error) {
	ud, err := b.repo.FindById(ctx, id)
	if ud == nil || err != nil {
		return nil, errcode.ErrorNotFound("用户不存在")
	}
	return ud, nil
}
