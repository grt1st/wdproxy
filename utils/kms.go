package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"strconv"
)

const (
	key     = "iyR6hu8Mn0iJPAKxQ3f9TCi7mFXXTRuI"
	keySize = 12
)

var aesBlock cipher.Block

func init() {
	var err error
	aesBlock, err = aes.NewCipher([]byte(key))
	if err != nil {
		panic(err.Error())
	}
}

func Encrypt(plaintext string) (string, error) {
	nonce := make([]byte, keySize)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	aesGcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return "", err
	}
	ciphertext := aesGcm.Seal(nil, nonce, []byte(plaintext), nil)
	ciphertext = append(nonce, ciphertext...)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func Decrypt(ciphertext string) (string, error) {
	crypted, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	nonce := crypted[0:keySize]
	cryBlob := crypted[keySize:]
	aesGcm, err := cipher.NewGCM(aesBlock)
	if err != nil {
		return "", err
	}
	var plaintext []byte
	plaintext, err = aesGcm.Open(nil, nonce, cryBlob, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func DecryptInt(ciphertext string) (int, error) {
	plaintext, err := Decrypt(ciphertext)
	if err == nil {
		return strconv.Atoi(plaintext)
	}
	return 0, err
}
