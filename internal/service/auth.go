package service

import (
	"context"
	v1 "github.com/coderlewin/kratosinit/api/proto/v1"
	"github.com/coderlewin/kratosinit/internal/biz"
	"github.com/coderlewin/kratosinit/internal/domain"
	"github.com/coderlewin/kratosinit/internal/pkg/ctxutils"
	"google.golang.org/protobuf/types/known/emptypb"
)

type AuthService struct {
	v1.UnimplementedAuthServer
	bz biz.IBiz
}

func NewAuthService(bz biz.IBiz) *AuthService {
	return &AuthService{bz: bz}
}

func (s *AuthService) Login(ctx context.Context, dto *v1.AuthLoginDTO) (*v1.AuthLoginVO, error) {
	token, err := s.bz.User().Login(ctx, dto.Account, dto.Password)
	if err != nil {
		return nil, err
	}
	return &v1.AuthLoginVO{AccessToken: token, TokenPrefix: "Bearer"}, nil
}

func (s *AuthService) Register(ctx context.Context, dto *v1.AuthRegisterDTO) (*emptypb.Empty, error) {
	err := s.bz.User().Register(ctx, dto.GetAccount(), dto.GetPassword(), dto.GetCheckPassword())
	return nil, err
}

func (s *AuthService) Logout(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	accessToken := ctxutils.FromAccessToken(ctx)
	err := s.bz.User().Logout(ctx, accessToken)
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *AuthService) LoginUserInfo(ctx context.Context, empty *emptypb.Empty) (*v1.UserVO, error) {
	id := ctxutils.FromUserID(ctx)
	user, err := s.bz.User().FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	return &v1.UserVO{
		Id:         user.ID,
		NickName:   user.NickName,
		Account:    user.Account,
		Avatar:     user.Avatar,
		UnionId:    user.UnionID,
		MpOpenId:   user.MpOpenID,
		Profile:    user.Profile,
		Role:       user.Role,
		CreateTime: user.CreateTime,
		UpdateTime: user.UpdateTime,
	}, nil
}

func (s *AuthService) UpdateMineInfo(ctx context.Context, dto *v1.UpdateMineInfoDTO) (*emptypb.Empty, error) {
	userId := ctxutils.FromUserID(ctx)
	err := s.bz.User().Update(ctx, &domain.User{
		ID:       userId,
		NickName: dto.GetNickName(),
		Avatar:   dto.GetAvatar(),
		Profile:  dto.GetProfile(),
	})
	return nil, err
}
