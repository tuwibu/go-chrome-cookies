package crypt

import (
	"crypto/aes"
	"crypto/cipher"
)

func ChromeDecrypt(key, encryptPass []byte) ([]byte, error) {
	return aesGCMDecrypt(encryptPass[15:], key, encryptPass[3:15])
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
