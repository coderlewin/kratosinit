package data

import (
	"context"
	"fmt"
	"github.com/bytedance/sonic"
	"github.com/coderlewin/kratosinit/internal/data/gorm_gen/dal"
	"github.com/coderlewin/kratosinit/internal/data/gorm_gen/entity"
	"github.com/coderlewin/kratosinit/internal/domain"
	"github.com/coderlewin/kratosinit/internal/pkg/constant"
	"github.com/coderlewin/kratosinit/internal/pkg/utils"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gen"
	"time"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

func NewUserRepo(data *Data, logger log.Logger) domain.UserRepo {
	return &userRepo{data: data, log: log.NewHelper(logger)}
}

func (u *userRepo) Find(ctx context.Context, conditions ...gen.Condition) (*domain.User, error) {
	qw := u.data.DB(ctx).User
	user, err := qw.WithContext(ctx).Where(conditions...).First()
	return u.convertToDomain(user), err
}

func (u *userRepo) FindByPage(ctx context.Context, current, size int, ud *domain.User) ([]*domain.User, int64, error) {
	qw := u.data.DB(ctx).User
	query := qw.WithContext(ctx)
	if ud.NickName != "" {
		query = query.Where(qw.NickName.Like("%" + ud.NickName + "%"))
	}
	if ud.Role != "" {
		query = query.Where(qw.Role.Eq(ud.Role))
	}
	list, count, err := query.FindByPage((current-1)*size, size)
	users := slice.Map(list, func(index int, item *entity.User) *domain.User {
		return u.convertToDomain(item)
	})
	return users, count, err
}

func (u *userRepo) Count(ctx context.Context, conditions ...gen.Condition) (int64, error) {
	return u.data.DB(ctx).User.WithContext(ctx).Where(conditions...).Count()
}

func (u *userRepo) FindById(ctx context.Context, id int64) (ud *domain.User, err error) {
	userJson, err := u.data.cache.Get(context.Background(), fmt.Sprintf(constant.CacheUserKey, id)).Result()
	if userJson != "" && err == nil {
		_ = sonic.UnmarshalString(userJson, &ud)
		return ud, nil
	}

	ud, err = u.Find(ctx, dal.User.ID.Eq(id))
	go func() {
		userStr, _ := sonic.MarshalString(ud)
		e := u.data.cache.Set(context.Background(), fmt.Sprintf(constant.CacheUserKey, id), userStr, time.Minute*5).Err()
		if e != nil {
			u.log.Errorf("set user cache error: %v", e)
		}
	}()
	return ud, err
}

func (u *userRepo) Create(ctx context.Context, user *domain.User) error {
	qw := u.data.DB(ctx).User
	timestamp := time.Now().UnixMilli()
	ue := u.convertToEntity(user)
	// id 使用雪花算法生成
	ue.ID = utils.GenSnowflakeID()
	ue.CreateTime = timestamp
	ue.UpdateTime = timestamp
	return qw.WithContext(ctx).Create(ue)
}

func (u *userRepo) Update(ctx context.Context, user *domain.User) error {
	qw := u.data.DB(ctx).User
	timestamp := time.Now().UnixMilli()
	user.UpdateTime = timestamp
	_, err := qw.WithContext(ctx).Where(qw.ID.Eq(user.ID)).Updates(u.convertToEntity(user))
	return err
}

func (u *userRepo) Delete(ctx context.Context, id int64) error {
	qw := u.data.DB(ctx).User
	_, err := qw.WithContext(ctx).Where(qw.ID.Eq(id)).Delete()
	return err
}

func (u *userRepo) FindByAccount(ctx context.Context, account string) (*domain.User, error) {
	return u.Find(ctx, dal.User.Account.Eq(account))
}

func (u *userRepo) convertToEntity(user *domain.User) *entity.User {
	if user == nil {
		return nil
	}
	return &entity.User{
		ID:         user.ID,
		Account:    user.Account,
		Password:   user.Password,
		UnionID:    user.UnionID,
		MpOpenID:   user.MpOpenID,
		NickName:   user.NickName,
		Avatar:     user.Avatar,
		Profile:    user.Profile,
		Role:       user.Role,
		CreateTime: user.CreateTime,
		UpdateTime: user.UpdateTime,
	}
}

func (u *userRepo) convertToDomain(user *entity.User) *domain.User {
	if user == nil {
		return nil
	}
	return &domain.User{
		ID:         user.ID,
		Account:    user.Account,
		Password:   user.Password,
		UnionID:    user.UnionID,
		MpOpenID:   user.MpOpenID,
		NickName:   user.NickName,
		Avatar:     user.Avatar,
		Profile:    user.Profile,
		Role:       user.Role,
		CreateTime: user.CreateTime,
		UpdateTime: user.UpdateTime,
	}
}
