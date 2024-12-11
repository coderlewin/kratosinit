package auth

import (
	"context"
	"fmt"
	"github.com/coderlewin/kratosinit/api/proto/errcode"
	"github.com/coderlewin/kratosinit/internal/conf"
	"github.com/duke-git/lancet/v2/convertor"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	ErrSignTokenFailed = errcode.ErrorUnauthorized("Failed to sign token")
	ErrTokenInvalid    = errcode.ErrorUnauthorized("Token is invalid")
	ErrTokenExpired    = errcode.ErrorUnauthorized("Token is expired")
)

const (
	ExpiredTokenTime = 2 * time.Hour
	cacheKeyPrefix   = "kratosinit:token:"
)

var ProviderSet = wire.NewSet(NewAuthnInterface)

type AuthnInterface interface {
	Sign(ctx context.Context, userId int64) (string, error)
	Destroy(ctx context.Context, accessToken string) error
	Verify(ctx context.Context, accessToken string) (int64, error)
}

type authnImpl struct {
	cache redis.Cmdable
	c     *conf.Jwt
}

func NewAuthnInterface(cache redis.Cmdable, c *conf.Jwt) AuthnInterface {
	return &authnImpl{cache: cache, c: c}
}

func (a *authnImpl) Sign(ctx context.Context, userId int64) (string, error) {
	now := time.Now()
	expiresAt := now.Add(ExpiredTokenTime)

	claims := jwt.NewWithClaims(jwt.SigningMethodHS512, &jwt.RegisteredClaims{
		Issuer:    "kratosinit",
		Subject:   convertor.ToString(userId),
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
	})
	token, err := claims.SignedString([]byte(a.c.GetSecret()))
	// 将 token 缓存
	err = a.cache.Set(ctx, fmt.Sprintf("%s%s", cacheKeyPrefix, token), "1", ExpiredTokenTime).Err()
	if err != nil {
		return "", ErrSignTokenFailed
	}
	return token, nil
}

func (a *authnImpl) Destroy(ctx context.Context, accessToken string) error {
	return a.cache.Del(ctx, fmt.Sprintf("%s%s", cacheKeyPrefix, accessToken)).Err()
}

func (a *authnImpl) Verify(ctx context.Context, accessToken string) (int64, error) {
	token, err := jwt.ParseWithClaims(accessToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(a.c.GetSecret()), nil
	})
	if err != nil || !token.Valid || token.Method != jwt.SigningMethodHS512 {
		return 0, ErrTokenInvalid
	}

	cmd := a.cache.Exists(ctx, fmt.Sprintf("%s%s", cacheKeyPrefix, accessToken))
	if err := cmd.Err(); err != nil {
		return 0, ErrTokenExpired
	}
	isExist := cmd.Val() > 0
	if !isExist {
		return 0, ErrTokenExpired
	}
	userId, _ := convertor.ToInt(token.Claims.(*jwt.RegisteredClaims).Subject)
	return userId, nil
}
