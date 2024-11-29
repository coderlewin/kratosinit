package server

import (
  v1 "github.com/coderlewin/kratosinit/api/proto/v1"
  "github.com/coderlewin/kratosinit/internal/conf"
  "github.com/coderlewin/kratosinit/internal/pkg/constant"
  middleware2 "github.com/coderlewin/kratosinit/internal/pkg/middleware"
  "github.com/coderlewin/kratosinit/internal/service"
  "github.com/go-kratos/kratos/v2/log"
  "github.com/go-kratos/kratos/v2/middleware/recovery"
  "github.com/go-kratos/kratos/v2/transport/http"
)

// NewHTTPServer new an HTTP server.
func NewHTTPServer(
  c *conf.Server,
  user *service.UserService,
  auth *service.AuthService,
  authMiddleware *middleware2.CheckAuthMiddleware,
  checkRole *middleware2.CheckRoleMiddleware,
  logger log.Logger) *http.Server {
  var opts = []http.ServerOption{
    http.Middleware(
      recovery.Recovery(),
      authMiddleware.Handle,
      checkRole.Handle(constant.RoleAdmin),
    ),
  }
  if c.Http.Network != "" {
    opts = append(opts, http.Network(c.Http.Network))
  }
  if c.Http.Addr != "" {
    opts = append(opts, http.Address(c.Http.Addr))
  }
  if c.Http.Timeout != nil {
    opts = append(opts, http.Timeout(c.Http.Timeout.AsDuration()))
  }
  srv := http.NewServer(opts...)
  v1.RegisterUserHTTPServer(srv, user)
  v1.RegisterAuthHTTPServer(srv, auth)
  return srv
}
