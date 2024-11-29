package user

import (
	"context"
	"github.com/coderlewin/kratosinit/api/proto/errcode"
	"github.com/coderlewin/kratosinit/internal/data/gorm_gen/dal"
	"github.com/coderlewin/kratosinit/internal/domain"
	"github.com/duke-git/lancet/v2/strutil"
)

func (b *Biz) Register(ctx context.Context, account, password, checkPassword string) error {
	if strutil.IsBlank(account) || strutil.IsBlank(password) || strutil.IsBlank(checkPassword) {
		return errcode.ErrorInvalidParameter("参数不合法")
	}

	if len(account) < 6 {
		return errcode.ErrorInvalidParameter("账号错误")
	}
	if len(password) < 8 {
		return errcode.ErrorInvalidParameter("密码错误")
	}
	if password != checkPassword {
		return errcode.ErrorInvalidParameter("两次密码输入不一致")
	}

	b.lock.Lock()
	defer b.lock.Unlock()
	count, err := b.repo.Count(ctx, dal.User.Account.Eq(account))
	if err != nil {
		b.log.Errorf("database error: query user failed: %v", err)
		return errcode.ErrorUnknown("系统内部异常")
	}
	if count > 0 {
		return errcode.ErrorInvalidParameter("账号已存在")
	}
	ud := &domain.User{Account: account}
	ud.EncryptPassword(password)
	err = b.repo.Create(ctx, ud)
	if err != nil {
		b.log.Errorf("database error: insert user failed: %v", err)
		return errcode.ErrorUnknown("系统内部异常")
	}
	return nil
}
