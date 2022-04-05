package utils

import (
	"encoding/json"
)

func StrToMap(s string) map[string]interface{} {
	m := make(map[string]interface{})
	json.Unmarshal([]byte(s), &m)
	return m
}

func MapToStr(m map[string]interface{}) string {
	bData, _ := json.Marshal(m)
	return string(bData)
}
