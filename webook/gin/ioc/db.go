package ioc

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"new_home/webook/gin/internal/repository/dao"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13306)/webook"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
