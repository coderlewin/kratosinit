package middleware

import (
	"context"
	v1 "github.com/coderlewin/kratosinit/api/proto/v1"
	"github.com/coderlewin/kratosinit/internal/conf"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"slices"
	"strings"
)

type CheckAuthMiddleware struct {
	c         *conf.Jwt
	whiteList []string
}

func NewCheckAuthMiddleware(c *conf.Jwt) *CheckAuthMiddleware {
	return &CheckAuthMiddleware{
		c: c,
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
		mapClaims := jwt.MapClaims{}
		parse, err := jwt.ParseWithClaims(token, &mapClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.c.Secret), nil
		})
		if err != nil || !parse.Valid {
			return nil, errors.New(http.StatusUnauthorized, "UNAUTHORIZED", "用户未授权")
		}
		userId, ok := mapClaims["userId"]
		if !ok {
			return nil, errors.New(http.StatusUnauthorized, "UNAUTHORIZED", "用户未授权")
		}

		ctx = context.WithValue(ctx, "userId", userId)
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
