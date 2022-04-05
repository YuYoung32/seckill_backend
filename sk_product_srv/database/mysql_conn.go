package database

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"sk_product_srv/conf"
)

var db *gorm.DB

func ConnInit() {
	dbConf := conf.GetDbConf()
	dsn := dbConf.Username + ":" + dbConf.Password + "@tcp(" + dbConf.Host + ":" + dbConf.Port + ")/" + dbConf.DbName + "?charset=utf8&parseTime=True&loc=Local"
	var err error
	db, err = gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}
	db.DB().SetMaxIdleConns(dbConf.MaxIdleConns)
	db.DB().SetMaxOpenConns(dbConf.MaxOpenConns)
	logrus.Info("mysql connect success")
}

func GetDBConn() *gorm.DB {
	return db
}
