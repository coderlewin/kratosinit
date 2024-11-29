package user

import (
	"github.com/coderlewin/kratosinit/internal/conf"
	"github.com/coderlewin/kratosinit/internal/domain"
	"github.com/go-kratos/kratos/v2/log"
	"sync"
)

type Biz struct {
	repo domain.UserRepo
	log  *log.Helper
	c    *conf.Jwt
	lock sync.Mutex
}

func NewBiz(repo domain.UserRepo, logger log.Logger, c *conf.Jwt) *Biz {
	return &Biz{
		repo: repo,
		log:  log.NewHelper(log.With(logger, "module", "biz/user")),
		c:    c,
	}
}
