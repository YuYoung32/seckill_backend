package database

import (
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"time"
)

type User struct {
	/*
		One "user" belongs to one "order“
	*/
	gorm.Model
	Username    string `json:"username" gorm:"type:varchar(100);not null"`
	Email       string `json:"email" gorm:"type:varchar(100);not null"`
	Password    string `json:"password" gorm:"type:varchar(200);not null"`
	Description string `json:"description" gorm:"type:varchar(100);not null"`
	Status      string `json:"status" gorm:"type:varchar(50);not null"`
}

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
		One "event" has many "order"
	*/
	gorm.Model
	Name        string    `gorm:"type:varchar(100);not null"`
	EventPrice  float32   `gorm:"type:float;not null"`
	EventNum    int       `gorm:"type:int;not null"`
	StartTime   time.Time `gorm:"type:datetime;not null"`
	EndTime     time.Time `gorm:"type:datetime;not null"`
	Description string    `gorm:"type:varchar(200)"`
	ProductId   uint      `gorm:"type:int"`

	Orders []Order `gorm:"foreignKey:EventID;references:ID"`
}

type Order struct {
	/*
		One "order" has one "event"
		One "order" has one "user"
	*/
	gorm.Model
	OrderSerial string `gorm:"type:varchar(100);not null"`
	PayStatus   string `gorm:"type:varchar(100);not null"` // 0:未支付 1:已支付

	EventId uint `gorm:"type:int;not null"`
	UserId  uint `gorm:"type:int;not null"`

	User User `gorm:"foreignKey:UserId;references:ID"`
}

//ModelInit 每次启动初始化数据库，便于更改数据库结构
func ModelInit() {
	//user表结构应在product服务里更改，不应在这里更改
	GetDBConn().AutoMigrate(&User{}, &Product{}, &Event{}, &Order{})
	logrus.Info("Auto migrate [User] [Product] [Event] [Order] success.")
}
