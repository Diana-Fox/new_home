package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

type UserDAO interface {
	FindByEmail(ctx context.Context, email string) (User, error)
	FindById(ctx context.Context, id int64) (User, error)
	Edit(ctx context.Context, id int64, u User) error
	Insert(ctx context.Context, user User) error
	FindByPhone(ctx context.Context, phone string) (User, error)
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &userDAO{
		db: db,
	}
}

type User struct {
	Id         int64          `gorm:"primaryKey,autoIncrement"`
	Email      sql.NullString `gorm:"unique"`
	Password   string
	Phone      sql.NullString `gorm:"unique"`
	NickName   string
	BrotherDay int64
	Biography  string
	//创建时间
	Ctime int64
	//更新时间
	Utime int64
}
type userDAO struct {
	db *gorm.DB
}

func (dao *userDAO) FindByPhone(ctx context.Context, phone string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("phone=?", phone).Find(&u).Error
	return u, err
}

func (dao *userDAO) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const uniqueConflictsErrNo = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			//邮箱冲突
			return ErrUserDuplicateEmail
		}
	}
	return nil
}

func (dao *userDAO) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).Find(&u).Error
	return u, err
}

func (dao *userDAO) FindById(ctx context.Context, id int64) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("id=?", id).Find(&u).Error
	return u, err
}

func (dao *userDAO) Edit(ctx context.Context, id int64, u User) error {
	u.Utime = time.Now().UnixMilli()
	err := dao.db.WithContext(ctx).Where("id=?", id).Updates(u).Error
	return err
}
