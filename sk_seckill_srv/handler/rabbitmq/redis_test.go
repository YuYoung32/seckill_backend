package rabbitmq

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"sk_seckill_srv/conf"
	"testing"
)

func TestRedisConn(t *testing.T) {
	conf.Init("../../conf/config.json")

	config := conf.GetRedisConf()

	conn, err := redis.Dial(config.Network, config.Host+":"+config.Port, redis.DialPassword(config.Password))
	if err != nil {
		t.Error(err)
		return
	}
	defer conn.Close()

	_, err = conn.Do("set", "mymap", map[string]interface{}{"name": "sk", "age": 18})
	if err != nil {
		t.Error(err)
		return
	}

	name, err := redis.String(conn.Do("get", "mymap"))
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(name)
	fmt.Printf("%T\n", name)

}
