package handler

import (
	"fmt"
	"sk_seckill_srv/conf"
	"sk_seckill_srv/database"
	"testing"
	"time"
)

func Test(t *testing.T) {
	conf.Init("../conf/config.json")
	database.ConnInit()
	database.ModelInit()
	db := database.GetDBConn()

	nowTime := time.Now()
	type ProductSelling struct {
		ProductId   int
		ProductName string
		LeftNum     int
	}
	var productSelling []ProductSelling
	res := db.Table("events").Select("events.product_id as product_id, products.name as product_name, products.left_num as left_num").
		Joins("left join products on events.product_id = products.id").
		Where("events.start_time <= ?", nowTime).Where("events.end_time >= ?", nowTime).
		Scan(&productSelling)
	fmt.Println(productSelling, res.RowsAffected)
}
