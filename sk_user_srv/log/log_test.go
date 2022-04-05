package log

import (
	"github.com/sirupsen/logrus"
	"sk_user_srv/conf"
	"testing"
)

func TestLog(t *testing.T) {
	conf.Init("../conf/config.json")
	logrus.Info("test")
}
