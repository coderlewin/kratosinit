package user

import (
	"context"
	"github.com/coderlewin/kratosinit/api/proto/errcode"
	"github.com/coderlewin/kratosinit/internal/domain"
)

func (b *Biz) FindByPage(ctx context.Context, current, size int, ud *domain.User) ([]*domain.User, int64, error) {
	users, total, err := b.repo.FindByPage(ctx, current, size, ud)
	if err != nil {
		b.log.Errorf("query user page error: %v", err)
		return nil, 0, errcode.ErrorUnknown("系统内部异常")
	}
	return users, total, nil
}
