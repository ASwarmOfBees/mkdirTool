package app

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"mkdirTool/public"
	"mkdirTool/tool"
)

//获取指定目录及所有子目录下的所有文件，可以匹配后缀过滤。
func WalkDir(dirPath string) error {
	var files []string

	dir, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return err
	}

	var oneMessage string
	for _, fi := range dir {
		if fi.IsDir() { // 开启新线程遍历目录
			absolutePath := dirPath + "/" + fi.Name()
			if public.IgnoreDir(fi.Name(), 2) == false {
				public.DirQueue.Push(absolutePath)
				public.WaitGroup.Wrap(func() {
					if dirName := public.DirQueue.Pop(); dirName != nil {
						dirStr := fmt.Sprintf("%s", dirName.Value)
						WalkDir(dirStr)
					}
				})
			} else {
				fmt.Println(fmt.Sprintf("directory \"%s\"ignore", absolutePath))
			}

			continue
		}

		//每一行一个文件，用逗号隔开，前面是文件名称，后面是哈希值，文件大小

		oneMessage = dirPath + "/" + fi.Name() + "," //文件名称
		if public.IgnoreFile(fi.Name(), 2) == true {
			fmt.Println(fmt.Sprintf("file \"%s\"ignore", oneMessage))
			continue
		}

		if hashValue, hashErr := tool.SHA1File(dirPath + "/" + fi.Name()); hashErr == nil {
			oneMessage += hashValue + "," //哈希值
		}

		oneMessage += strconv.FormatInt(fi.Size(), 10) //文件大小

		files = append(files, oneMessage)
	}
	if len(files) != 0 {
		messageStr := tool.String2Bytes(strings.Join(files, "\n"))

		_, _ = tool.WriteFile(public.Config.SaveFile, messageStr)
		fmt.Println(fmt.Sprintf("directory \"%s\" ergodic completed", dirPath))
	}

	return nil
}

//获取指定目录及所有子目录下的所有文件。
func WalkRecursionDir(dirPath string) (files []string, err error) {
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
		hashValue, hashErr := tool.SHA1File(filename)

		if hashErr == nil {
			oneMessage += hashValue + "," //哈希值
		}

		oneMessage += strconv.FormatInt(file.Size(), 10) //文件大小

		messageStr := tool.String2Bytes(oneMessage + "\n")

		_, _ = tool.WriteFile(public.Config.SaveFile, messageStr)

		//	fmt.Println(messageStr)

		files = append(files, oneMessage)

		return nil
	})

	return files, err
}
