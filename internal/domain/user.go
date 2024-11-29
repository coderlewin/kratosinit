package domain

import (
	"context"
	"github.com/duke-git/lancet/v2/cryptor"
	"gorm.io/gen"
)

type UserRepo interface {
	FindByAccount(ctx context.Context, account string) (*User, error)
	FindById(ctx context.Context, id int64) (*User, error)
	Create(ctx context.Context, user *User) error
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id int64) error
	Count(ctx context.Context, conditions ...gen.Condition) (int64, error)
	FindByPage(ctx context.Context, current, size int, ud *User) ([]*User, int64, error)
	Find(ctx context.Context, conditions ...gen.Condition) (*User, error)
}

type User struct {
	ID         int64  // id
	Account    string // 账号
	Password   string // 密码
	UnionID    string // 微信开放平台id
	MpOpenID   string // 公众号openId
	NickName   string // 用户昵称
	Avatar     string // 用户头像
	Profile    string // 用户简介
	Role       string // 用户角色：user/admin/ban
	CreateTime int64  // 创建时间
	UpdateTime int64  // 更新时间
}

func (u *User) CheckPassword(pwd string) bool {
	return cryptor.Md5String(pwd) == u.Password
}

func (u *User) EncryptPassword(pwd string) {
	// md5 加密
	u.Password = cryptor.Md5String(pwd)
}
