package predis

import (
	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
	"sk_seckill_srv/conf"
	"time"
)

var redisClient *redis.Pool

func Init() {
	config := conf.GetRedisConf()
	redisClient = &redis.Pool{
		MaxIdle:     config.MaxIdle,
		MaxActive:   config.MaxActive,
		IdleTimeout: time.Second * time.Duration(config.MaxIdleTimeout),
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial(config.Network, config.Host+":"+config.Port, redis.DialPassword(config.Password))
			if err != nil {
				logrus.WithField("module", "redis_init").Error(err)
				return nil, err
			}
			return con, nil
		},
	}
}

func GetRedisConn() redis.Conn {
	return redisClient.Get()
}
