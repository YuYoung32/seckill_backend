package predis

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"sk_web/conf"
	"testing"
)

func TestGetRedisConn(t *testing.T) {
	conf.Init("../conf/config.json")
	Init()

	conn := GetRedisConn()
	defer conn.Close()
	res, err := redis.String(conn.Do("GET", "testuser0@qq.com"))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(res)
}
