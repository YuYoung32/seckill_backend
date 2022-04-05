/*
	配置类型，与conf.json中的配置类型严格对应
*/

package conf

type ConfigType struct {
	DbConf       `json:"database"`
	LogConf      `json:"logrus"`
	RabbitmqConf `json:"rabbitmq"`
	RedisConf    `json:"redis"`
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

type RabbitmqConf struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Vhost    string `json:"vhost"`
}

type RedisConf struct {
	Network        string `json:"network"`
	Host           string `json:"host"`
	Port           string `json:"port"`
	Password       string `json:"password"`
	MaxIdle        int    `json:"max_idle"`
	MaxActive      int    `json:"max_active"`
	MaxIdleTimeout int    `json:"max_idle_timeout"`
}
