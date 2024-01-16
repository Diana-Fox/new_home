package repository

import (
	"context"
	"database/sql"
	"new_home/webook/gin/internal/domain"
	"new_home/webook/gin/internal/repository/cache"
	"new_home/webook/gin/internal/repository/dao"
	"time"
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
	FindByPhone(ctx context.Context, phone string) (domain.User, error)
}
type userRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func (r *userRepository) FindByPhone(ctx context.Context, phone string) (domain.User, error) {
	u, err := r.dao.FindByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{Id: u.Id, Email: u.Email.String,
		PassWord: u.Password}, nil
}

func NewUserRepository(dao dao.UserDAO, cache cache.UserCache) UserRepository {
	return &userRepository{
		dao:   dao,
		cache: cache,
	}
}

func (r *userRepository) Create(ctx context.Context, u domain.User) error {
	return r.dao.Insert(ctx, r.domainToEntity(u))
}
func (r *userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := r.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return r.entityToDomain(u), nil
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
	u = r.entityToDomain(user)
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
func (r *userRepository) domainToEntity(u domain.User) dao.User {
	return dao.User{
		Id: u.Id,
		Email: sql.NullString{String: u.Email,
			Valid: u.Email != ""},
		Password: u.PassWord,
		Phone: sql.NullString{
			String: u.Phone,
			Valid:  u.Phone != "",
		},
		Ctime: u.Ctime.UnixMilli(),
	}
}
func (r *userRepository) entityToDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email.String,
		PassWord: u.Password,
		Phone:    u.Phone.String,
		Ctime:    time.UnixMilli(u.Ctime),
	}
}
