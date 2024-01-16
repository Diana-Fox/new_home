package domain

import (
	"time"
)

// 领域对象
type User struct {
	Id         int64
	Email      string
	PassWord   string
	NickName   string
	BrotherDay int64
	Biography  string
	Ctime      time.Time
	Phone      string
}
