package main

import (
	"os"
	"runtime"
	"strings"
)

func SaveStaticFile(fileName, content string) error {
	path := ""
	if runtime.GOOS == "windows" {
		path = ".\\static\\"
	} else {
		path = "./static/"
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		if errCreatingDir := os.Mkdir(path, os.ModeAppend); err != nil {
			return errCreatingDir
		}
		file, errCreatingFile := os.Create(path + fileName)
		if err != nil {
			return errCreatingFile
		}
		defer file.Close()
		_, errWritingInFile := file.WriteString(content)
		if errWritingInFile != nil {
			return errWritingInFile
		}
		return nil
	}
	return os.ErrExist
}

func GetRootUrlByFindingThirdIndexOfCharacter(str string, char string) string {
	urlCopy := str
	indexGlobal := 0
	for i := 0; i < 3; i++ {
		index := strings.Index(urlCopy, char)
		if index == -1 {
			// ИСХОДИМ ИЗ ТОГО, ЧТО ЕСЛИ -1, ТО В КОНЦЕ НЕ ПОСТАВЛЕНА / -> URL ТИПА http://example.com
			return str + "/"
		}
		indexGlobal = indexGlobal + index + 1
		urlCopy = urlCopy[index+1:]
	}
	str = str[:indexGlobal]
	return str
}

func GetDirByFindingThirdIndexOfCharacter(str string, char string) string {
	urlCopy := str
	indexGlobal := 0
	for i := 0; i < 3; i++ {
		index := strings.Index(urlCopy, char)
		if index == -1 {
			// ИСХОДИМ ИЗ ТОГО, ЧТО ЕСЛИ -1, ТО В КОНЦЕ НЕ ПОСТАВЛЕНА / -> URL ТИПА http://example.com
			return "/"
		}
		indexGlobal = indexGlobal + index + 1
		urlCopy = urlCopy[index+1:]
	}
	str = str[indexGlobal-1:]
	return str
}

func getLastPartOfDir(dir string) string {
	lastIndexOfSlash := strings.LastIndex(dir, "/")
	return dir[lastIndexOfSlash:]
}

func getDirWithoutLastPath(dir string) string {
	lastIndexOfSlash := strings.LastIndex(dir, "/")
	if lastIndexOfSlash == 0 {
		return ""
	}
	if !strings.HasSuffix(dir, "/") {
		return dir[:lastIndexOfSlash] + "/"
	}
	return dir[:lastIndexOfSlash]
}

func GetFileInUrl(str string, char string) string {
	rootUrl := ""
	if strings.HasPrefix(str, "http") {
		rootUrl = GetRootUrlByFindingThirdIndexOfCharacter(str, "/")
		if rootUrl == str+"/" {
			// ЗНАЧИТ ПОЛУЧИЛИ СТРОКУ ТИПА http://example.com
			return ""
		}
		return rootUrl
	} else {
		lastIndex := strings.LastIndex(str, char)
		if lastIndex == -1 {
			return ""
		}
		fileName := str[lastIndex+1:]
		return fileName
	}
}

func ValidateFormatFile(fileName string) bool {
	formats := []string{"css", "ico", "js", "woff", "woff2", "ttf", "eot"}
	for _, value := range formats {
		if strings.HasSuffix(fileName, "."+value) {
			return true
		}
	}
	return false
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func removeDuplicates(strList []string) []string {
	list := []string{}
	for _, item := range strList {
		if contains(list, item) == false {
			list = append(list, item)
		}
	}
	return list
}
