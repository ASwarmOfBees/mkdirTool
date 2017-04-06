package tool

import (
	"fmt"
	"io/ioutil"
	"os"
)

//读取文件内容
func ReadFile(filePath string) ([]byte, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}

func IsFileEmpty(filePath string) bool {
	stat, err := os.Stat(filePath)
	var empty = true
	if err == nil {
		if stat.Size() == 0 {
			empty = true
		} else {
			empty = false
		}
	}

	return empty
}

func CheckFileIsExist(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

//写文件
func WriteFile(filename string, bytes []byte) (int, error) {
	if len(filename) == 0 {
		fmt.Println("config file is lost")

		return -1, nil
	}

	var f *os.File
	var err error

	if CheckFileIsExist(filename) { //如果文件存在
		if f, err = os.OpenFile(filename, os.O_APPEND, 0666); err == nil { //打开文件
			//添加换行符
			if fInfo, errStat := f.Stat(); errStat == nil {
				if fInfo.Size() != 0 {
					f.Write([]byte("\r\n"))
				}
			}
		}

	} else {
		f, err = os.Create(filename) //创建文件
		fmt.Println("create file success")
	}

	if err != nil {
		fmt.Println("config file open fail")
		return -1, err
	}

	n, result := f.Write(bytes)
	f.Close()

	return n, result
}
