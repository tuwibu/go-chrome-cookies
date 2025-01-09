package main

import (
	"encoding/json"
	"fmt"

	"github.com/tuwibu/go-chrome-cookies/data"
	"github.com/tuwibu/go-chrome-cookies/logger"
)

func main() {
	logger.InitLogger()
	config := data.NewConfig("C:\\Users\\zorovhs\\.multiprofile\\user-data-dir\\532153507730427904-copy")
	cookies, err := config.LoadCookies()
	if err != nil {
		fmt.Println(err)
	}
	jsonData, err := json.Marshal(cookies)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(jsonData))
}
