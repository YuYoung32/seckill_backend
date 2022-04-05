package conf

import (
	"encoding/json"
	"io/ioutil"
)

var GlobalConf ConfigType

func Init(filePath string) {
	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(buf, &GlobalConf)
	if err != nil {
		panic(err)
	}
}

func GetDbConf() *DbConf {
	return &GlobalConf.DbConf
}

func GetLogConf() *LogConf {
	return &GlobalConf.LogConf
}

func GetRabbitmqConf() *RabbitmqConf {
	return &GlobalConf.RabbitmqConf
}

func GetRedisConf() *RedisConf {
	return &GlobalConf.RedisConf
}
