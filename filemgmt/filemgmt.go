package filemgmt

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/tuwibu/go-chrome-cookies/logger"
)

const Prefix = "[go-cc]: "

func IntToBool(a int) bool {
	switch a {
	case 0, -1:
		return false
	}
	return true
}

func TimeStampFormat(stamp int64) time.Time {
	s1 := time.Unix(stamp, 0)
	if s1.Local().Year() > 9999 {
		return time.Date(9999, 12, 13, 23, 59, 59, 0, time.Local)
	}
	return s1
}

func TimeEpochFormat(epoch int64) time.Time {
	maxTime := int64(99633311740000000)
	if epoch > maxTime {
		return time.Date(2049, 1, 1, 1, 1, 1, 1, time.Local)
	}

	// Chuyển đổi từ Windows FILETIME (100-nanosecond intervals since January 1, 1601 UTC)
	// sang Unix epoch (seconds since January 1, 1970 UTC)
	epochMicros := (epoch - 116444736000000000) / 10 // Convert to microseconds
	seconds := epochMicros / 1000000
	micros := epochMicros % 1000000

	return time.Unix(seconds, micros*1000)
}

func ReadFile(filename string) (string, error) {
	s, err := ioutil.ReadFile(filename)
	return string(s), err
}

func WriteFile(filename string, data []byte) error {
	err := ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return nil
	}
	return err
}

func FormatFileName(dir, browser, filename, format string) string {
	r := strings.Replace(strings.TrimSpace(strings.ToLower(browser)), " ", "_", -1)
	p := path.Join(dir, fmt.Sprintf("%s_%s.%s", r, filename, format))
	return p
}

func MakeDir(dirName string) error {
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		return os.Mkdir(dirName, 0700)
	}
	return nil
}

func Compress(exportDir string) error {
	files, err := ioutil.ReadDir(exportDir)
	if err != nil {
		logger.NewLogger().Error(err)
	}
	var b = new(bytes.Buffer)
	zw := zip.NewWriter(b)
	for _, f := range files {
		fw, _ := zw.Create(f.Name())
		fileName := path.Join(exportDir, f.Name())
		fileContent, err := ioutil.ReadFile(fileName)
		if err != nil {
			_ = zw.Close()
			return err
		}
		_, err = fw.Write(fileContent)
		if err != nil {
			_ = zw.Close()
			return err
		}
		err = os.Remove(fileName)
		if err != nil {
			logger.NewLogger().Error(err)
		}
	}
	if err := zw.Close(); err != nil {
		return err
	}
	zipName := exportDir + `/archive.zip`
	outFile, _ := os.Create(zipName)
	_, err = b.WriteTo(outFile)
	if err != nil {
		return err
	}
	fmt.Printf("%s Compress success, zip filename is %s \n", Prefix, zipName)
	return nil
}

func CloseFile() func(*os.File) {
	return func(f *os.File) {
		_ = f.Close()
	}
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func TimeToWindowsEpoch(t time.Time) int64 {
	// Chuyển từ Unix time sang Windows FILETIME
	epochMicros := t.UnixMicro()
	return (epochMicros * 10) + 116444736000000000
}
