package user

import (
  "github.com/coderlewin/kratosinit/internal/domain"
  "github.com/coderlewin/kratosinit/internal/pkg/auth"
  "github.com/go-kratos/kratos/v2/log"
  "sync"
)

type Biz struct {
  repo  domain.UserRepo
  log   *log.Helper
  lock  sync.Mutex
  authn auth.AuthnInterface
}

func NewBiz(repo domain.UserRepo, logger log.Logger, authn auth.AuthnInterface) *Biz {
  return &Biz{
    repo:  repo,
    log:   log.NewHelper(log.With(logger, "module", "biz/user")),
    authn: authn,
  }
}
