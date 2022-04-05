package utils

import (
	"encoding/base64"
	"io/ioutil"
)

//Pic2Base64 converts a picture to base64 string
func Pic2Base64(picPath string) (string, error) {
	pic, err := ioutil.ReadFile(picPath)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(pic), nil
}
