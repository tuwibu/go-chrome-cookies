package crypt

import (
	"crypto/aes"
	"crypto/cipher"

	"github.com/tuwibu/go-chrome-cookies/throw"
)

func ChromeDecrypt(key, encryptPass []byte) ([]byte, error) {
	if len(encryptPass) > 15 {
		return aesGCMDecrypt(encryptPass[15:], key, encryptPass[3:15])
	} else {
		return nil, throw.ErrorPasswordIsEmpty()
	}
}

func aesGCMDecrypt(encrypted, key, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockMode, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	origData, err := blockMode.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return nil, err
	}
	return origData, nil
}
