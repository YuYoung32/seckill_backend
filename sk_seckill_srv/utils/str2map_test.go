package utils

import (
	"fmt"
	"testing"
)

func TestStrToMap(t *testing.T) {
	var str = `{"key1":"value1","key2":"value2"}`
	var m = StrToMap(str)
	fmt.Println(m)
}
