package utils

import (
	"fmt"
	"testing"
)

func TestEncryption(t *testing.T) {
	hash, _ := Encryption("123456")
	fmt.Println(hash)
	fmt.Println(len(hash))
}

func TestValidationPassword(t *testing.T) {
	hash, _ := Encryption("123456")
	fmt.Println(ValidationPassword(hash, "123456"))
}
