package crypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
)

func ChromeEncrypt(key, plaintext []byte) ([]byte, error) {
	// Chrome sử dụng AES-GCM với:
	// - 3 byte prefix "v10"
	// - 12 byte nonce
	// - Ciphertext
	// - 16 byte auth tag

	// Mỗi lần encrypt sẽ tạo nonce mới ngẫu nhiên
	nonce := make([]byte, 12)
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	// Tạo cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Tạo GCM mode
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Mã hóa
	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)

	// Tạo kết quả cuối cùng theo format của Chrome:
	// v10 + nonce + ciphertext
	result := make([]byte, 3+len(nonce)+len(ciphertext))
	copy(result[0:3], []byte("v10"))
	copy(result[3:15], nonce)
	copy(result[15:], ciphertext)

	return result, nil
}

func aesGCMEncrypt(plaintext, key, nonce []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	return blockMode.Seal(nil, nonce, plaintext, nil), nil
}
