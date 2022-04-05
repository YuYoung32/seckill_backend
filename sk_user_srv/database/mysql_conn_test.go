package database

import (
	"sk_user_srv/conf"
	"testing"
)

func TestGetDatabaseConn(t *testing.T) {
	conf.Init("../conf/config.json")
	ConnInit()
	conn := GetDBConn()
	if conn == nil {
		t.Error("Database connection is nil")
	}
}
