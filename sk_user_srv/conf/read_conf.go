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

func GetEmailConf() *EmailConf {
	return &GlobalConf.EmailConf
}

func GetCacheConf() *CacheConf {
	return &GlobalConf.CacheConf
}
