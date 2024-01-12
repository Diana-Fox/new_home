package repository

import (
	"context"
	"new_home/webook/gin/internal/domain"
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
	dao dao.UserDAO
}

func NewUserRepository(dao dao.UserDAO) UserRepository {
	return &userRepository{
		dao: dao,
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
	u, err := r.dao.FindById(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	return domain.User{Id: u.Id,
		Email:      u.Email,
		NickName:   u.NickName,
		BrotherDay: u.BrotherDay,
		Biography:  u.Biography}, nil
}

func (r *userRepository) Edit(ctx context.Context, u domain.User) error {
	return r.dao.Edit(ctx, u.Id, dao.User{
		NickName:   u.NickName,
		Biography:  u.Biography,
		BrotherDay: u.BrotherDay,
	})
}
