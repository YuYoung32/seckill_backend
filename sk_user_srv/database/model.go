package database

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type User struct {
	gorm.Model
	Username    string `json:"username" gorm:"type:varchar(100);not null"`
	Email       string `json:"email" gorm:"type:varchar(100);not null"`
	Password    string `json:"password" gorm:"type:varchar(200);not null"`
	Description string `json:"description" gorm:"type:varchar(100);not null"`
	Status      string `json:"status" gorm:"type:varchar(50);not null"`
}

type Admin struct {
	gorm.Model
	Username string `json:"username" gorm:"type:varchar(100);not null"`
	Password string `json:"password" gorm:"type:varchar(200);not null"`
}

//ModelInit 每次启动初始化数据库，便于更改数据库结构
func ModelInit() {
	GetDBConn().AutoMigrate(&User{}, &Admin{})
	logrus.Info("Auto migrate [User] [Admin] success.")
}
