package util

import (
	"crypto/md5"
	"io"
	"os"
)

func MD5Sum(filePath string) ([]byte, error) {
	in, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer in.Close()
	m5h := md5.New()
	_, err = io.Copy(m5h, in)
	if err != nil {
		return nil, err
	}
	return m5h.Sum([]byte(``)), nil
}
