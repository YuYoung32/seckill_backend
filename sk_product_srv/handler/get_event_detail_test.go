package handler

import (
	"fmt"
	"sk_product_srv/conf"
	"sk_product_srv/database"
	"testing"
	"time"
)

type EventDetailedInfo_ struct {
	EventId         string
	EventCreateTime time.Time
	Name            string
	EventPrice      float32
	EventNum        int
	StartTime       time.Time
	EndTime         time.Time
	ProductId       uint
	ProductName     string
}

func TestGetDetailEvent(t *testing.T) {
	conf.Init("../conf/config.json")
	database.ConnInit()
	db := database.GetDBConn()

	var res []EventDetailedInfo_
	db.Debug().Table("events").Select("events.id as event_id, events.created_at as event_create_time, events.name, " +
		"events.event_price, events.event_num, events.start_time ,events.end_time, events.product_id, products.name as product_name").
		Joins("left join products on events.product_id = products.id").Offset(0).Limit(10).
		Scan(&res)

	fmt.Println(res)
}

func TestGetDetailEvents(t *testing.T) {
	conf.Init("../conf/config.json")
	database.ConnInit()
	db := database.GetDBConn()
	var event EventDetailedInfo_
	db.Debug().Table("events").Where("events.id=?", "1").Select("events.id as event_id, events.created_at as event_create_time, events.name, events.event_price, events.event_num, events.start_time ,events.end_time, events.product_id, products.name as product_name").
		Joins("left join products on events.product_id = products.id").Scan(&event)
	fmt.Println(event)
}

func TestProductImpl_GetSelectedProductList(t *testing.T) {
	conf.Init("../conf/config.json")
	database.ConnInit()
	db := database.GetDBConn()

	var products []database.Product
	db.Debug().Where("id !=?", 1).Find(&products)
}

func TestEventImpl_GetFrontEventList(t *testing.T) {
	conf.Init("../conf/config.json")
	database.ConnInit()
	db := database.GetDBConn()

	var events []EventDetailedWithProduct
	db.Table("events").Joins("left join products on events.product_id = products.id").
		Select("events.id            as event_id," +
			"       events.event_num     as left_num," +
			"       events.event_price   as event_price," +
			"       events.start_time    as start_time," +
			"       events.end_time      as end_time," +
			"       products.id          as product_id," +
			"       products.name        as product_name," +
			"       products.picture     as picture," +
			"       products.price       as price," +
			"       products.description as product_description," +
			"       products.unit        as unit").
		Scan(&events)

	fmt.Println(events)
}
