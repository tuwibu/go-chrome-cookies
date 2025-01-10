package utils

import (
	"io"
	"os"
)

func ReadFile(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func WriteFile(path string, content string) error {
	return os.WriteFile(path, []byte(content), 0644)
}

func IsFileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func CheckOrCreateDir(path string) error {
	return os.MkdirAll(path, 0755)
}

func CopyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	buf := make([]byte, 32*1024)

	for {
		n, err := sourceFile.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		if _, err := destinationFile.Write(buf[:n]); err != nil {
			return err
		}
	}
	return nil
}
