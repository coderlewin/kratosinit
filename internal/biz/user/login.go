package user

import (
	"context"
	"github.com/coderlewin/kratosinit/api/proto/errcode"
	"github.com/duke-git/lancet/v2/strutil"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func (b *Biz) Login(ctx context.Context, account, password string) (string, error) {
	if strutil.IsBlank(account) || strutil.IsBlank(password) {
		return "", errcode.ErrorInvalidParameter("账号或密码不能为空")
	}
	if len(account) < 6 {
		return "", errcode.ErrorInvalidParameter("账号错误")
	}
	if len(password) < 8 {
		return "", errcode.ErrorInvalidParameter("密码错误")
	}

	ud, err := b.repo.FindByAccount(ctx, account)
	if ud == nil || err != nil {
		return "", errcode.ErrorNotFound("账号不存在")
	}
	isValid := ud.CheckPassword(password)
	if !isValid {
		return "", errcode.ErrorInvalidParameter("账号与密码不匹配")
	}
	// 生成 token
	unix := time.Now().Unix()
	s := time.Second * 60 * 60
	token, err := b.getJwtToken(b.c.Secret, unix, int64(s.Seconds()), ud.ID)
	if err != nil {
		b.log.Error("jwt token generate error: ", err.Error())
		return "", errcode.ErrorUnknown("jwt token 生成失败")
	}
	return token, nil
}

func (b *Biz) getJwtToken(secret string, iat, seconds int64, userId int64) (string, error) {
	claims := make(jwt.MapClaims)
	claims["exp"] = iat + seconds
	claims["iat"] = iat
	claims["userId"] = userId
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims = claims
	return token.SignedString([]byte(secret))
}
