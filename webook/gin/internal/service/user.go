package service

import (
	"context"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"new_home/webook/gin/internal/domain"
	"new_home/webook/gin/internal/repository"
	"time"
)

var (
	ErrUserDuplicateEmail    = repository.ErrUserDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("账号或密码不对")
	ErrUserNotFound          = errors.New("用户不存在")
	ErrTimeFormat            = errors.New("日期转化失败")
)

type UserService interface {
	SignUp(ctx context.Context, u domain.User) error
	Login(ctx context.Context, email, password string) (domain.User, error)
	Edit(ctx context.Context, id int64, nickName string, dataStr string, biography string) error
	Profile(ctx context.Context, id int64) (domain.User, error)
	FindOrCreate(ctx context.Context, phone string) (domain.User, error)
}
type userService struct {
	repo repository.UserRepository
}

func (svc *userService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	//先查再说
	user, err := svc.repo.FindByPhone(ctx, phone)
	if err != repository.ErrUserNotFind {
		//查到或者报错不是未查到
		return user, err
	}
	u := domain.User{Phone: phone}
	//去注册
	err = svc.repo.Create(ctx, u)
	if err != nil {
		return domain.User{}, err
	}
	//再查一次，但是这里会有主从延迟问题
	return svc.repo.FindByPhone(ctx, phone)
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

// 加context是为了保持链路
func (svc *userService) SignUp(ctx context.Context, u domain.User) error {
	//加密和存储
	password, err := bcrypt.GenerateFromPassword([]byte(u.PassWord), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.PassWord = string(password)
	return svc.repo.Create(ctx, u)
}
func (svc *userService) Login(ctx context.Context, email, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err == repository.ErrUserNotFind {
		return u, ErrUserDuplicateEmail
	}
	if err != nil {
		return u, err
	}
	encrypted, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	err = bcrypt.CompareHashAndPassword([]byte(encrypted), []byte(password))
	if err != nil {
		//没通过呗,debug或者info日志
		return u, ErrInvalidUserOrPassword
	}

	return u, nil
}

func (svc *userService) Edit(ctx context.Context, id int64, nickName string, dataStr string, biography string) error {
	u, err := svc.repo.FindById(ctx, id)
	if err != nil {
		return err
	}
	if u.Id == 0 {
		return ErrUserNotFound
	}
	time, err := time.Parse("2006-01-02", dataStr)
	if err != nil {
		return ErrTimeFormat
	}
	u.BrotherDay = time.UnixMilli()
	u.NickName = nickName
	u.Biography = biography
	return svc.repo.Edit(ctx, u)
}

func (svc *userService) Profile(ctx context.Context, id int64) (domain.User, error) {
	u, err := svc.repo.FindById(ctx, id)
	if u.Id == 0 {
		return u, ErrUserNotFound
	}
	return u, err
}
