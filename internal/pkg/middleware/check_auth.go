package middleware

import (
	"context"
	v1 "github.com/coderlewin/kratosinit/api/proto/v1"
	"github.com/coderlewin/kratosinit/internal/pkg/auth"
	"github.com/coderlewin/kratosinit/internal/pkg/ctxutils"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"slices"
	"strings"
)

type CheckAuthMiddleware struct {
	whiteList []string
	authn     auth.AuthnInterface
}

func NewCheckAuthMiddleware(authn auth.AuthnInterface) *CheckAuthMiddleware {
	return &CheckAuthMiddleware{
		authn: authn,
		whiteList: []string{
			v1.OperationAuthLogin,
			v1.OperationAuthRegister,
		},
	}
}

func (m *CheckAuthMiddleware) Handle(h middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		tr, _ := transport.FromServerContext(ctx)
		operation := tr.Operation()
		if slices.Contains(m.whiteList, operation) {
			return h(ctx, req)
		}
		token := m.extractTokenString(tr)
		userId, err := m.authn.Verify(context.Background(), token)
		if err != nil {
			return nil, err
		}
		ctx = ctxutils.NewUserID(ctx, userId)
		ctx = ctxutils.NewAccessToken(ctx, token)
		return h(ctx, req)
	}
}

func (m *CheckAuthMiddleware) extractTokenString(tr transport.Transporter) string {
	authCode := tr.RequestHeader().Get("Authorization")
	if authCode == "" {
		return ""
	}
	// SplitN 的意思是切割字符串，但是最多 N 段
	// 如果要是 N 为 0 或者负数，则是另外的含义，可以看它的文档
	authSegments := strings.SplitN(authCode, " ", 2)
	if len(authSegments) != 2 {
		// 格式不对
		return ""
	}
	return authSegments[1]
}
