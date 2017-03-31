package tool

import (
	"crypto/sha1"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

//生成指定文件sha1值
func SHA1File(path string) (string, error) {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return "", err
	}

	h := sha1.New()
	_, err = io.Copy(h, file)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

//获取指定目录及所有子目录下的所有文件，可以匹配后缀过滤。
func WalkDir(dirPath string) (files []string, err error) {
	files = make([]string, 0, 30)
	var oneMessage string

	err = filepath.Walk(dirPath, func(filename string, file os.FileInfo, err error) error { //遍历目录
		if err != nil { //忽略错误
			return err
		}

		if file.IsDir() { // 忽略目录
			return nil
		}

		oneMessage = filename + "," //文件名称

		//每一行一个文件，用逗号隔开，前面是文件名称，后面是哈希值，文件大小
		hashValue, hashErr := SHA1File(filename)

		if hashErr == nil {
			oneMessage += hashValue + "," //哈希值
		}

		oneMessage += strconv.FormatInt(file.Size(), 10) //文件大小

		fmt.Println(oneMessage)

		files = append(files, oneMessage)

		return nil
	})

	return files, err
}
