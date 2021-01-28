package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// 判断所给路径文件/文件夹是否存在
func Exists(path string) bool {
	_, err := os.Stat(path) //os.Stat获取文件信息
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

// 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// 判断所给路径是否为文件
func IsFile(path string) bool {
	return !IsDir(path)
}

func GetFileList(root string) (files []string) {
	if err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	}); err != nil {
		log.Println("获取文件夹中列表失败", err)
	}
	return
}

//获取文件目录和文件名称（去掉后缀）
func GetFileName(filename string, suf string) (dir, suffix string) {
	dir = filepath.Dir(filename)
	base := filepath.Base(filename)
	suffix = strings.TrimSuffix(base, suf)
	return
}

//获取文件夹下的所有文件，支持嵌套的
func GetAllFiles(dirpath string) (files []string) {
	if Exists(dirpath) {
		List := GetFileList(dirpath)
		for _, f := range List {
			if !IsDir(f) {
				file := f
				files = append(files, file)
			}
		}
	}
	return
}

//输入一个文件名称，为其所在的目录创建文件夹
func MakeDirForFile(filepath string) error {
	dir, _ := GetFileName(filepath, "")
	err := os.MkdirAll(dir, 0766)
	if err != nil {
		e := fmt.Sprintf("创建文件夹%v失败: %v", dir, err)
		return errors.New(e)
	}
	return nil
}
