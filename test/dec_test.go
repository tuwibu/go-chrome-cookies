package test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/tuwibu/go-chrome-cookies/data"
	"github.com/tuwibu/go-chrome-cookies/logger"
)

func TestGetKey(t *testing.T) {
	logger.InitLogger()
	config := data.NewConfig("C:\\Users\\zorovhs\\.multiprofile\\user-data-dir\\532153507797536768")
	cookies, err := config.LoadCookies()
	if err != nil {
		fmt.Println(err)
	}
	jsonData, err := json.Marshal(cookies)
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile("a.json", jsonData, 0644)
	if err != nil {
		fmt.Println(err)
	}

}
