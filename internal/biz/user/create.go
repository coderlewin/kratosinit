package user

import (
	"context"
	"github.com/coderlewin/kratosinit/api/proto/errcode"
	"github.com/coderlewin/kratosinit/internal/domain"
	"github.com/duke-git/lancet/v2/strutil"
)

func (b *Biz) Create(ctx context.Context, ud *domain.User) error {
	if strutil.IsBlank(ud.Account) || strutil.IsBlank(ud.NickName) {
		return errcode.ErrorInvalidParameter("参数不合法")
	}
	if len(ud.Account) < 6 {
		return errcode.ErrorInvalidParameter("账号错误")
	}
	// 设置默认密码
	ud.EncryptPassword("12345678")
	// 加锁，防止重复创建
	b.lock.Lock()
	defer b.lock.Unlock()
	err := b.repo.Create(ctx, ud)
	if err != nil {
		b.log.Errorf("db create user error: %v", err)
		return errcode.ErrorOperationFailed("创建失败")
	}
	return nil
}
