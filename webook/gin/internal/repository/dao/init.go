package dao

import "gorm.io/gorm"

// 自动建表的一些操作。一般开发用不到
func InitTable(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
