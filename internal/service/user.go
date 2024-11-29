package service

import (
	"context"
	v1 "github.com/coderlewin/kratosinit/api/proto/v1"
	"github.com/coderlewin/kratosinit/internal/biz"
	"github.com/coderlewin/kratosinit/internal/domain"
	"github.com/duke-git/lancet/v2/slice"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserService struct {
	v1.UnimplementedUserServer
	bz biz.IBiz
}

func NewUserService(bz biz.IBiz) *UserService {
	return &UserService{bz: bz}
}

func (s *UserService) FindById(ctx context.Context, request *v1.IdRequest) (*v1.UserVO, error) {
	user, err := s.bz.User().FindById(ctx, request.Id)
	if err != nil {
		return nil, err
	}
	return s.convertToUserVO(user), nil
}

func (s *UserService) Create(ctx context.Context, dto *v1.UserCreateDTO) (*emptypb.Empty, error) {
	err := s.bz.User().Create(ctx, &domain.User{
		Account:  dto.GetAccount(),
		NickName: dto.GetNickName(),
		Avatar:   dto.GetAvatar(),
		Role:     dto.GetRole(),
	})
	return nil, err
}

func (s *UserService) Delete(ctx context.Context, request *v1.IdRequest) (*emptypb.Empty, error) {
	err := s.bz.User().Delete(ctx, request.Id)
	return nil, err
}

func (s *UserService) Update(ctx context.Context, dto *v1.UserUpdateDTO) (*emptypb.Empty, error) {
	err := s.bz.User().Update(ctx, &domain.User{
		ID:       dto.GetId(),
		NickName: dto.GetNickName(),
		Avatar:   dto.GetAvatar(),
		Profile:  dto.GetProfile(),
		Role:     dto.GetRole(),
	})
	return nil, err
}

func (s *UserService) PageList(ctx context.Context, dto *v1.UserQueryDTO) (*v1.UserPageVO, error) {
	list, total, err := s.bz.User().FindByPage(
		ctx,
		int(dto.GetCurrent()),
		int(dto.GetSize()),
		&domain.User{
			NickName: dto.GetNickName(),
			Role:     dto.GetRole(),
		})
	if err != nil {
		return nil, err
	}
	users := slice.Map(list, func(index int, item *domain.User) *v1.UserVO {
		return s.convertToUserVO(item)
	})
	return &v1.UserPageVO{
		Total: total,
		List:  users,
	}, nil
}

func (s *UserService) convertToUserVO(user *domain.User) *v1.UserVO {
	if user == nil {
		return nil
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
	}
}
