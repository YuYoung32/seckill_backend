package log

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"sk_seckill_srv/conf"
)

func Init() {
	logConf := conf.GetLogConf()

	level := map[string]logrus.Level{
		"debug": logrus.DebugLevel,
		"info":  logrus.InfoLevel,
		"warn":  logrus.WarnLevel,
		"error": logrus.ErrorLevel,
		"fatal": logrus.FatalLevel,
		"panic": logrus.PanicLevel,
	}
	logrus.SetLevel(level[logConf.Level])

	logrus.SetFormatter(&logrus.JSONFormatter{})

	file, err := os.OpenFile(logConf.FilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		panic(err)
	}

	logrus.SetOutput(io.MultiWriter(file, os.Stdout))

	logrus.Info("Logrus init success")
}
