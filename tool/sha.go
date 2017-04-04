package tool

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
)

//生成指定文件sha1值
func SHA1File(path string) (string, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return fmt.Sprintf("open %s file  fail", path), err
	}

	h := sha1.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return fmt.Sprintf("construct sha1 fail", path), err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
