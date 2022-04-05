/*
	配置类型，与conf.json中的配置类型严格对应
*/

package conf

type ConfigType struct {
	DbConf    `json:"database"`
	LogConf   `json:"logrus"`
	CacheConf `json:"cache"`
}

type DbConf struct {
	Host         string `json:"host"`
	Port         string `json:"port"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	DbName       string `json:"db_name"`
	MaxIdleConns int    `json:"max_idle_conns"`
	MaxOpenConns int    `json:"max_open_conns"`
}

type LogConf struct {
	Level    string `json:"level"`
	FilePath string `json:"file_path"`
}

type CacheConf struct {
	ExpireTime int `json:"expire_time"`
	CleanTime  int `json:"clean_time"`
}
