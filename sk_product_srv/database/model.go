package database

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"time"
)

type Product struct {
	/*
		One "product" has many "event"
	*/
	gorm.Model
	Name        string  `gorm:"type:varchar(100);not null"`
	Price       float32 `gorm:"type:float;not null"`
	LeftNum     int     `gorm:"type:int;not null"`
	Unit        string  `gorm:"type:varchar(50);not null"`
	Picture     string  `gorm:"type:varchar(100)"`
	Description string  `gorm:"type:varchar(200)"`
	Events      []Event `gorm:"foreignKey:ProductID;references:ID;constraint:OnDelete:SET NULL"`
}

type Event struct {
	/*
		One "event" has one "product"
	*/
	gorm.Model
	Name        string    `gorm:"type:varchar(100);not null"`
	EventPrice  float32   `gorm:"type:float;not null"`
	EventNum    int       `gorm:"type:int;not null"`
	StartTime   time.Time `gorm:"type:datetime;not null"`
	EndTime     time.Time `gorm:"type:datetime;not null"`
	Description string    `gorm:"type:varchar(200)"`
	ProductId   uint      `gorm:"type:int"`
}

//ModelInit 每次启动初始化数据库，便于更改数据库结构
func ModelInit() {
	GetDBConn().AutoMigrate(&Product{}, &Event{})
	logrus.Info("Auto migrate [Product Event] success.")
}
