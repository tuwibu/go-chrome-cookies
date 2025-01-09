package data

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/tuwibu/go-chrome-cookies/crypt"
	"github.com/tuwibu/go-chrome-cookies/filemgmt"
	"github.com/tuwibu/go-chrome-cookies/logger"
)

type cookie struct {
	Host         string
	Path         string
	KeyName      string
	encryptValue []byte
	Value        string
	IsSecure     bool
	IsHTTPOnly   bool
	HasExpire    bool
	IsPersistent bool
	CreateDate   time.Time
	ExpireDate   time.Time
}

func (c *Config) LoadCookies() (map[string][]cookie, error) {
	fmt.Println(c.CookiePath)
	cookieDB, err := sql.Open("sqlite3", c.CookiePath)
	if err != nil {
		return nil, err
	}
	defer cookieDB.Close()
	rows, err := cookieDB.Query("SELECT name, encrypted_value, host_key, path, creation_utc, expires_utc, is_secure, is_httponly, has_expires, is_persistent FROM cookies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			key, host, path                               string
			isSecure, isHTTPOnly, hasExpire, isPersistent int
			createDate, expireDate                        int64
			value, encryptValue                           []byte
		)
		err = rows.Scan(&key, &encryptValue, &host, &path, &createDate, &expireDate, &isSecure, &isHTTPOnly, &hasExpire, &isPersistent)
		if err != nil {
			logger.NewLogger().Error(err)
		}
		cookie := cookie{
			KeyName:      key,
			Host:         host,
			Path:         path,
			encryptValue: encryptValue,
			IsSecure:     filemgmt.IntToBool(isSecure),
			IsHTTPOnly:   filemgmt.IntToBool(isHTTPOnly),
			HasExpire:    filemgmt.IntToBool(hasExpire),
			IsPersistent: filemgmt.IntToBool(isPersistent),
			CreateDate:   filemgmt.TimeEpochFormat(createDate),
			ExpireDate:   filemgmt.TimeEpochFormat(expireDate),
		}
		// remove 'v10'
		if c.Key == "" {
			value, err = crypt.DPApi(encryptValue)
		} else {
			value, err = crypt.ChromePass([]byte(c.Key), encryptValue)
		}
		if err != nil {
			logger.NewLogger().Error(err)
		}
		cookie.Value = string(value)
		c.cookies[host] = append(c.cookies[host], cookie)
	}
	return c.cookies, nil
	// fmt.Println(c.Key)
	// return nil, nil
}

func (c *Config) SaveCookies(cookies map[string]string) error {
	return nil
}
