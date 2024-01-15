package repository

import (
	"context"
	"new_home/webook/gin/internal/domain"
	"new_home/webook/gin/internal/repository/cache"
	"new_home/webook/gin/internal/repository/dao"
)

var (
	ErrUserDuplicateEmail = dao.ErrUserDuplicateEmail
	ErrUserNotFind        = dao.ErrUserNotFound
)

type UserRepository interface {
	Create(ctx context.Context, u domain.User) error
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	FindById(ctx context.Context, id int64) (domain.User, error)
	Edit(ctx context.Context, u domain.User) error
}
type userRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDAO, cache cache.UserCache) UserRepository {
	return &userRepository{
		dao:   dao,
		cache: cache,
	}
}

func (r *userRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.PassWord,
	})
}
func (r *userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{Id: u.Id, Email: u.Email,
		PassWord: u.Password}, nil
}

func (r *userRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	u, err := r.cache.Get(ctx, id) //先去缓存找
	if err == nil {
		return u, err
	}
	if err == cache.ErrKeyNotExist {
		//查库
	}
	user, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	u = domain.User{Id: u.Id,
		Email:      user.Email,
		NickName:   user.NickName,
		BrotherDay: user.BrotherDay,
		Biography:  user.Biography}
	err = r.cache.Set(ctx, u)
	return u, nil
}

func (r *userRepository) Edit(ctx context.Context, u domain.User) error {
	return r.dao.Edit(ctx, u.Id, dao.User{
		NickName:   u.NickName,
		Biography:  u.Biography,
		BrotherDay: u.BrotherDay,
	})
}
