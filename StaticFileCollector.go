package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"strings"
)

/*
1. Получаем html +
2. Парсим html, получаем все ссылки
3. Получаем rootUrl +
4. Делаем запрос в каждый статический файл и получаем Body, приводим в строку
5. Сохраняем (в бд или в файлы)
6. Тесты
*/
func main1() {
	//	TestsServiceFunctions()
	url := "http://192.168.146.130/"
	rootUrl := GetRootUrlByFindingThirdIndexOfCharacter(url, "/")
	resp := MakeRequest(url, "/")
	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return
	}
	var hrefAndSrc []string
	doc.Find("*").Each(func(i int, selection *goquery.Selection) {
		href, hrefExists := selection.Attr("href")
		src, srcExists := selection.Attr("src")
		var dir string
		if hrefExists {
			dir = href
		}
		if srcExists {
			dir = src
		}
		if hrefExists || srcExists {
			if strings.Contains(dir, "://") {
				return
			}
			response := MakeRequest(rootUrl, dir)
			if response != nil {
				fileName, Content := GetFileNameAndHtmlFromResponse(response)
				if fileName != "" {
					// ПОМЕНЯТЬ НА ЗАПИСЬ В БД, А НЕ В СРЕЗ
					hrefAndSrc = append(hrefAndSrc, fileName)
					// УБРАТЬ. БУДЕТ ЗАПИСЬ В БД
					err := SaveStaticFile(fileName, Content)
					if err != nil {
						return
					}
				}
			}
		}
	})
	hrefAndSrc = removeDuplicates(hrefAndSrc)
	for _, value := range hrefAndSrc {
		fmt.Println(value)
	}
}

func MakeRequest(url string, dir string) *http.Response {
	if strings.HasSuffix(url, "/") && strings.HasPrefix(dir, "/") {
		// Doc: url = http://example.com/ && dir = /styles.css
		dir = dir[1:]
	} else {
		// Doc: url = http://example.com && dir = styles.css
		dir = "/" + dir
	}
	// Doc: url = http://example.com/ && dir = styles.css
	// Doc: url = http://example.com && dir = /styles.css
	// Doc: ТАКИЕ СИТУАЦИИ УДОВЛЕТВОРЯЮТ
	resp, err := http.Get(url + dir)
	if err != nil {
		return nil
	}
	return resp
}

func GetFileNameAndHtmlFromResponse(response *http.Response) (string, string) {
	path := GetFileInUrl(response.Request.URL.Path, "/")
	if !ValidateFormatFile(path) {
		path = ""
	}
	index := strings.LastIndex(path, "/")
	fileName := ""
	if index != -1 {
		fileName = path[index+1:]
	} else {
		fileName = strings.ReplaceAll(response.Request.URL.Path, "//", "/")
	}
	contentBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return fileName, "Error in reading Bytes"
	}
	contentString := string(contentBytes)
	return fileName, contentString
}

func TestsServiceFunctions() {
	url := "http://192.168.146.130/salavat/styles.css"
	url1 := "http://192.168.146.130/1/2/3/4/5/6/7/8/styles.css"
	url2 := "http://192.168.146.130/styles.css"
	url3 := "http://192.168.146.130/"
	url4 := "http://192.168.146.130"
	dir := GetDirByFindingThirdIndexOfCharacter(url, "/")
	dir1 := GetDirByFindingThirdIndexOfCharacter(url1, "/")
	dir2 := GetDirByFindingThirdIndexOfCharacter(url2, "/")
	dir3 := GetDirByFindingThirdIndexOfCharacter(url3, "/")
	dir4 := GetDirByFindingThirdIndexOfCharacter(url4, "/")
	if dir != "/salavat/styles.css" &&
		dir1 != "/1/2/3/4/5/6/7/8/styles.css" &&
		dir2 != "/styles.css" &&
		dir3 != "/" &&
		dir4 != "/" {
		panic("bad test GetDirByFindingThirdIndexOfCharacter")
	}
	rootUrl := GetRootUrlByFindingThirdIndexOfCharacter(url, "/")
	rootUrl1 := GetRootUrlByFindingThirdIndexOfCharacter(url1, "/")
	rootUrl2 := GetRootUrlByFindingThirdIndexOfCharacter(url2, "/")
	rootUrl3 := GetRootUrlByFindingThirdIndexOfCharacter(url3, "/")
	rootUrl4 := GetRootUrlByFindingThirdIndexOfCharacter(url4, "/")
	if rootUrl != "http://192.168.146.130/" &&
		rootUrl1 != "http://192.168.146.130/" &&
		rootUrl2 != "http://192.168.146.130/" &&
		rootUrl3 != "http://192.168.146.130/" &&
		rootUrl4 != "http://192.168.146.130/" {
		panic("bad test GetDirByFindingThirdIndexOfCharacter")
	}
	fileName := GetFileInUrl(url, "/")
	fileName1 := GetFileInUrl(url1, "/")
	fileName2 := GetFileInUrl(url2, "/")
	fileName3 := GetFileInUrl(url3, "/")
	fileName4 := GetFileInUrl(url4, "/")
	if fileName != "styles.css" &&
		fileName1 != "styles.css" &&
		fileName2 != "styles.css" &&
		fileName3 != "" &&
		fileName4 != "" {
		panic("bad test GetFileInUrl")
	}
	err := SaveStaticFile("styles.css", "bad state")
	if err != nil {
		panic("bad test SaveStaticFile")
	}
}
