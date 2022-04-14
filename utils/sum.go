package utils

import (
	"crypto/sha256"
	"fmt"
)

func Sha256(data []byte) string {
	sum := sha256.Sum256(data)
	return fmt.Sprintf("%x", sum)
}
