package paths

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// 判断所给路径文件/文件夹是否存在
func IsExist(path string) bool {
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

func CreateDir(path string) bool {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func WalkDir(path string) []string {
	fmt.Println(path)
	rd, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println(err)
		log.Println("theDirPath cannot walk")
		return []string{}
	}
	fileNames := []string{}
	for _, fi := range rd {
		if fi.IsDir() {
			fileNames = append(fileNames, fi.Name())
		}
	}
	return fileNames
}

func GetExeDir() string {
	exeDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return exeDir
}

func GetAbsPath(path string) string {
	absDir, err := filepath.Abs(path)
	if err != nil {
		log.Fatal(err)
	}
	return absDir
}