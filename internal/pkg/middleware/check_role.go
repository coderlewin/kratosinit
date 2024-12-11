package middleware

import (
	"context"
	"github.com/coderlewin/kratosinit/api/proto/errcode"
	"github.com/coderlewin/kratosinit/internal/domain"
	"slices"
	"time"

	v1 "github.com/coderlewin/kratosinit/api/proto/v1"
	"github.com/coderlewin/kratosinit/internal/pkg/ctxutils"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
)

type CheckRoleMiddleware struct {
	list     []string
	userRepo domain.UserRepo
}

func NewCheckRoleMiddleware(repo domain.UserRepo) *CheckRoleMiddleware {
	return &CheckRoleMiddleware{
		userRepo: repo,
		list: []string{
			v1.OperationUserCreate,
			v1.OperationUserDelete,
			v1.OperationUserPageList,
			v1.OperationUserFindById,
			v1.OperationUserUpdate,
		},
	}
}

func (m *CheckRoleMiddleware) Handle(role string) func(h middleware.Handler) middleware.Handler {
	return func(h middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			tr, _ := transport.FromServerContext(ctx)
			operation := tr.Operation()
			if !slices.Contains(m.list, operation) {
				return h(ctx, req)
			}

			userId := ctxutils.FromUserID(ctx)
			// 创建一个带超时的新的 context
			newCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			ud, err := m.userRepo.FindById(newCtx, userId)
			if ud == nil || err != nil {
				return nil, errcode.ErrorUnauthorized("用户未授权")
			}
			if ud.Role != role {
				return nil, errcode.ErrorForbidden("用户权限不足")
			}
			return h(ctx, req)

		}
	}
}
