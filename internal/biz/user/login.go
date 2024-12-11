package user

import (
  "context"
  "github.com/coderlewin/kratosinit/api/proto/errcode"
  "github.com/duke-git/lancet/v2/strutil"
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
  token, err := b.authn.Sign(ctx, ud.ID)
  if err != nil {
    b.log.Error("jwt token generate error: ", err.Error())
    return "", err
  }
  return token, nil
}
