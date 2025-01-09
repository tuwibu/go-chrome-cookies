package data

import (
	"encoding/base64"
	"path/filepath"
	"strings"

	"github.com/tuwibu/go-chrome-cookies/crypt"
	"github.com/tuwibu/go-chrome-cookies/utils"
)

func getKey(folder string) string {
	localStatePath := filepath.Join(folder, "Local State")
	localState, err := utils.ReadFile(localStatePath)
	if err != nil {
		return ""
	}
	if !strings.Contains(localState, "encrypted_key") {
		return ""
	}
	encryptedKey := strings.Split(localState, "encrypted_key\":\"")[1]
	encryptedKey = strings.Split(encryptedKey, "\"")[0]
	encryptedKey = strings.TrimSpace(encryptedKey)
	key, err := base64.StdEncoding.DecodeString(encryptedKey)
	if err != nil {
		return ""
	}
	decryptedKey, err := crypt.DPApi(key[5:])
	if err != nil {
		return ""
	}
	return string(decryptedKey)
}
