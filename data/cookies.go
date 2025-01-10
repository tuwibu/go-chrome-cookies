package data

import (
	"database/sql"
	"embed"
	"errors"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/tuwibu/go-chrome-cookies/crypt"
	"github.com/tuwibu/go-chrome-cookies/filemgmt"
	"github.com/tuwibu/go-chrome-cookies/logger"
	"github.com/tuwibu/go-chrome-cookies/utils"
)

type Cookie struct {
	Host                 string
	Path                 string
	KeyName              string
	encryptValue         []byte
	Value                string
	IsSecure             bool
	IsHTTPOnly           bool
	HasExpire            bool
	IsPersistent         bool
	CreateDate           time.Time
	ExpireDate           time.Time
	LastAccessDate       time.Time
	LastUpdateDate       time.Time
	TopFrameSiteKey      string
	Priority             int
	SameSite             int
	SourceScheme         int
	SourcePort           int
	SourceType           int
	HasCrossSiteAncestor bool
}

//go:embed Cookies
var Cookies embed.FS

func (c *Config) LoadCookies() (map[string][]Cookie, error) {
	if !utils.IsFileExists(c.CookiePath) {
		return nil, errors.New("cookie file not found")
	}
	cookieDB, err := sql.Open("sqlite3", c.CookiePath)
	if err != nil {
		return nil, err
	}
	defer cookieDB.Close()
	rows, err := cookieDB.Query("SELECT name, encrypted_value, host_key, path, creation_utc, expires_utc, is_secure, is_httponly, has_expires, is_persistent, last_access_utc, last_update_utc, top_frame_site_key, priority, samesite, source_scheme, source_port, source_type, has_cross_site_ancestor FROM cookies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var (
			key, host, path, topFrameSiteKey                                                                                              string
			isSecure, isHTTPOnly, hasExpire, isPersistent, priority, sameSite, sourceScheme, sourcePort, sourceType, hasCrossSiteAncestor int
			createDate, expireDate, lastAccessDate, lastUpdateDate                                                                        int64
			value, encryptValue                                                                                                           []byte
		)
		err = rows.Scan(&key, &encryptValue, &host, &path, &createDate, &expireDate, &isSecure, &isHTTPOnly, &hasExpire, &isPersistent, &lastAccessDate, &lastUpdateDate, &topFrameSiteKey, &priority, &sameSite, &sourceScheme, &sourcePort, &sourceType, &hasCrossSiteAncestor)
		if err != nil {
			logger.NewLogger().Error(err)
		}
		cookie := Cookie{
			KeyName:         key,
			Host:            host,
			Path:            path,
			TopFrameSiteKey: topFrameSiteKey,
			Priority:        priority,
			SameSite:        sameSite,
			SourceScheme:    sourceScheme,
			SourcePort:      sourcePort,
			encryptValue:    encryptValue,
			IsSecure:        filemgmt.IntToBool(isSecure),
			IsHTTPOnly:      filemgmt.IntToBool(isHTTPOnly),
			HasExpire:       filemgmt.IntToBool(hasExpire),
			IsPersistent:    filemgmt.IntToBool(isPersistent),
			CreateDate:      filemgmt.TimeEpochFormat(createDate),
			ExpireDate:      filemgmt.TimeEpochFormat(expireDate),
			LastAccessDate:  filemgmt.TimeEpochFormat(lastAccessDate),
			LastUpdateDate:  filemgmt.TimeEpochFormat(lastUpdateDate),
		}
		value, err = crypt.ChromeDecrypt([]byte(c.Key), encryptValue)
		if err != nil {
			logger.NewLogger().Error(err)
		}
		cookie.Value = string(value)
		c.cookies[host] = append(c.cookies[host], cookie)
	}
	return c.cookies, nil
}

func (c *Config) AddCookie(cookie Cookie) {
	// check if cookie.Host is in c.cookies
	if _, ok := c.cookies[cookie.Host]; !ok {
		c.cookies[cookie.Host] = []Cookie{}
	}
	// check if cookie.KeyName is in c.cookies[cookie.Host], rewrite if exist
	for i, existingCookie := range c.cookies[cookie.Host] {
		if existingCookie.KeyName == cookie.KeyName {
			existingCookie.Value = cookie.Value
			existingCookie.IsSecure = cookie.IsSecure
			existingCookie.IsHTTPOnly = cookie.IsHTTPOnly
			existingCookie.HasExpire = cookie.HasExpire
			existingCookie.IsPersistent = cookie.IsPersistent
			existingCookie.CreateDate = cookie.CreateDate
			existingCookie.ExpireDate = cookie.ExpireDate
			existingCookie.LastAccessDate = cookie.LastAccessDate
			existingCookie.LastUpdateDate = cookie.LastUpdateDate
			existingCookie.TopFrameSiteKey = cookie.TopFrameSiteKey
			existingCookie.Priority = cookie.Priority
			existingCookie.SameSite = cookie.SameSite
			existingCookie.SourceScheme = cookie.SourceScheme
			existingCookie.SourcePort = cookie.SourcePort
			existingCookie.SourceType = cookie.SourceType
			existingCookie.HasCrossSiteAncestor = cookie.HasCrossSiteAncestor
			c.cookies[cookie.Host][i] = existingCookie
			return
		}
	}
	c.cookies[cookie.Host] = append(c.cookies[cookie.Host], cookie)
}

func (c *Config) SaveCookies(cookies map[string][]Cookie) error {
	if !utils.IsFileExists(c.CookiePath) {
		_ = c.InitCookies()
	}
	// copy file to .bk
	cookieDB, err := sql.Open("sqlite3", c.CookiePath)
	if err != nil {
		return err
	}
	defer cookieDB.Close()

	// Bắt đầu transaction
	tx, err := cookieDB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT OR REPLACE INTO cookies 
		(host_key, name, value, encrypted_value, path, creation_utc, expires_utc, 
		is_secure, is_httponly, has_expires, is_persistent, last_access_utc, last_update_utc, 
		top_frame_site_key, priority, samesite, source_scheme, source_port, source_type, 
		has_cross_site_ancestor) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, cookieList := range cookies {
		for _, cookie := range cookieList {
			// Validate các trường quan trọng
			if cookie.Host == "" || cookie.KeyName == "" || cookie.Path == "" {
				continue // Skip invalid cookies
			}

			// Validate và normalize các giá trị số
			priority := normalizeInt(cookie.Priority, 1)      // Default: 1
			sameSite := normalizeInt(cookie.SameSite, -1)     // Default: -1
			sourcePort := normalizeInt(cookie.SourcePort, -1) // Default: -1

			// Mã hóa giá trị cookie
			encryptValue, err := crypt.ChromeEncrypt([]byte(c.Key), []byte(cookie.Value))
			if err != nil {
				logger.NewLogger().Error(err)
				continue
			}

			// Chuyển đổi thời gian sang Windows FILETIME format
			createTime := filemgmt.TimeToWindowsEpoch(cookie.CreateDate)
			expireTime := filemgmt.TimeToWindowsEpoch(cookie.ExpireDate)
			lastAccessTime := filemgmt.TimeToWindowsEpoch(cookie.LastAccessDate)
			lastUpdateTime := filemgmt.TimeToWindowsEpoch(cookie.LastUpdateDate)

			_, err = stmt.Exec(
				cookie.Host,
				cookie.KeyName,
				"", // Chrome luôn để trống trường value
				encryptValue,
				cookie.Path,
				createTime,
				expireTime,
				filemgmt.BoolToInt(cookie.IsSecure),
				filemgmt.BoolToInt(cookie.IsHTTPOnly),
				filemgmt.BoolToInt(cookie.HasExpire),
				filemgmt.BoolToInt(cookie.IsPersistent),
				lastAccessTime,
				lastUpdateTime,
				cookie.TopFrameSiteKey,
				priority,
				sameSite,
				cookie.SourceScheme,
				sourcePort,
				cookie.SourceType,
				filemgmt.BoolToInt(cookie.HasCrossSiteAncestor),
			)
			if err != nil {
				logger.NewLogger().Error(err)
				return err
			}
		}
	}

	// Commit transaction
	return tx.Commit()
}

func (c *Config) InitCookies() error {
	content, err := Cookies.ReadFile("Cookies")
	if err != nil {
		return err
	}

	return utils.WriteFile(c.CookiePath, string(content))
}

func normalizeInt(value, defaultValue int) int {
	if value == 0 {
		return defaultValue
	}
	return value
}
