package test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/tuwibu/go-chrome-cookies/data"
)

func TestEncrypt(t *testing.T) {
	config := data.NewConfig("C:\\Users\\zorovhs\\.multiprofile\\user-data-dir\\532153507730427904-copy")
	jsonData, err := os.ReadFile("a.json")
	if err != nil {
		fmt.Println(err)
	}
	cookies := make(map[string][]data.Cookie)
	err = json.Unmarshal(jsonData, &cookies)
	if err != nil {
		fmt.Println(err)
	}
	if err := config.SaveCookies(cookies); err != nil {
		fmt.Println(err)
	}
	fmt.Println("Saved")
}
