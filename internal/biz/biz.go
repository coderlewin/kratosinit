package biz

import (
	"context"
	"github.com/coderlewin/kratosinit/internal/biz/user"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(New, user.NewBiz)

type Transaction interface {
	InTx(context.Context, func(ctx context.Context) error) error
}

type IBiz interface {
	User() *user.Biz
}

type biz struct {
	userBiz *user.Biz
}

func New(ubiz *user.Biz) IBiz {
	return &biz{userBiz: ubiz}
}

func (b *biz) User() *user.Biz {
	return b.userBiz
}
