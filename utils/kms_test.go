package utils

import (
	"fmt"
	"testing"
)

func TestEncrypt(t *testing.T) {
	data, err := Encrypt("abcdefg")
	if err != nil {
		t.Fatalf(err.Error())
	}
	plain, err := Decrypt(data)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(plain == "abcdefg")
}
