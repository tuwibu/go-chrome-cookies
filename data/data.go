package data

import "path/filepath"

type Config struct {
	Folder     string
	Key        string
	CookiePath string
	cookies    map[string][]cookie
}

func NewConfig(folder string) *Config {

	return &Config{
		Folder:     folder,
		CookiePath: filepath.Join(folder, "Default", "Network", "Cookies"),
		Key:        getKey(folder),
		cookies:    make(map[string][]cookie),
	}
}

func (c *Config) GetKey() string {
	return c.Key
}
