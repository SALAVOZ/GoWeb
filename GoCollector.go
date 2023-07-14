package main

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/css/scanner"
	"io"
	"net/http"
	"strings"
)

type Page struct {
	path    string
	content string
}

type pagesSlice []Page

var pagesAll pagesSlice
var baseRootUrl string

func main() {
	url := "/../1/fonts.www"
	b, a, err := strings.Cut(url, "..")
	fmt.Println(b+a, err)
	dump("http://192.168.146.130/")
}

func dump(url string) {
	var pagesTMP pagesSlice
	pagesAll = pagesTMP
	rootUrl := GetRootUrlByFindingThirdIndexOfCharacter(url, "/")
	baseRootUrl = rootUrl
	// ЗАПИСЬ В срез pages
	pages := dumpSiteRecursively(rootUrl, "/")
	fmt.Println(len(pages))
}

func dumpSite(rootUrl, dir string) pagesSlice {
	var pagesCurrent pagesSlice
	response := doRequest(rootUrl, dir)
	if response.StatusCode != 200 {
		return pagesCurrent
	}
	doc, err := goquery.NewDocumentFromResponse(response)
	if err != nil {
		return pagesCurrent
	}
	doc.Find("*").Each(func(i int, selection *goquery.Selection) {
		var dir string
		href, hrefExists := selection.Attr("href")
		src, srcExists := selection.Attr("src")
		if hrefExists {
			dir = href
		}
		if srcExists {
			dir = src
		}
		if !strings.HasPrefix(dir, "/") && dir != "/" && dir != "" {
			dir = "/" + dir
		}
		if ValidatePath(dir) && !strings.Contains(dir, "://") && !strings.Contains(dir, "#") {
			response = doRequest(rootUrl, dir)
			if response.StatusCode != 200 {
				return
			}
			contentBytes, err := io.ReadAll(response.Body)
			if err != nil {
				fmt.Println(err)
				return
			}
			path := GetDirByFindingThirdIndexOfCharacter(rootUrl, "/") + dir
			if strings.HasPrefix(path, "//") {
				path = path[1:]
			}
			contentString := string(contentBytes)
			pagesAll = append(pagesAll, Page{path, contentString})
			pagesCurrent = append(pagesCurrent, Page{dir, contentString})
		}
		return
	})
	return pagesCurrent
}

func dumpPage(page Page) {
	name := page.path
	content := page.content
	err := SaveStaticFile(name, content)
	if err != nil {
		fmt.Println(err)
	}
}

func doRequest(url, dir string) *http.Response {
	// Doc: url = http://example.com/ && dir = styles.css
	// Doc: url = http://example.com && dir = /styles.css
	// Doc: ТАКИЕ СИТУАЦИИ УДОВЛЕТВОРЯЮТ
	resp, err := http.Get(concatenateUrlAndDir(url, dir))
	if err != nil {
		return nil
	}
	return resp
}

func dumpSiteRecursively(rootUrl, dir string) pagesSlice {
	var pagesCurrent pagesSlice
	pagesCurrent = dumpSite(rootUrl, dir)
	if len(pagesCurrent) != 0 {
		for _, page := range pagesCurrent {
			if format, isValid := isFileStatic(page.path); isValid {
				switch format {
				case ".css":
					{
						urls := getAllSrcInCSSFile(page.content)
						for _, url := range urls {
							response := doRequest(baseRootUrl, url)
							bytesContent, err := io.ReadAll(response.Body)
							if err != nil {
								fmt.Println(err)
							}
							stringContent := string(bytesContent)
							pagesAll = append(pagesAll, Page{path: url, content: stringContent})
						}
					}
				case ".js":
					{
						fmt.Println(page.path)
					}
				case ".woff":
					{
						fmt.Println(page.path)
					}
				case ".woff2":
					{
						fmt.Println(page.path)
					}
				case ".ico":
					{
						fmt.Println(page.path)
					}
				}
			} else if strings.Contains(page.content, "<head") || strings.Contains(page.content, "<body") {
				//pagesAll = append(pagesCurrent, dumpSiteRecursively(rootUrl+getDirWithoutLastPath(page.path), getLastPartOfDir(page.path))...)
				dumpSiteRecursively(rootUrl+getDirWithoutLastPath(page.path), getLastPartOfDir(page.path))
			}
		}
	}
	return pagesCurrent
}

func isFileStatic(path string) (string, bool) {
	formats := []string{"css", "ico", "js", "woff", "woff2", "ttf", "eot", "svg"}
	for _, format := range formats {
		if strings.HasSuffix(path, "."+format) {
			return "." + format, true
		}
	}
	return "", false
}

func ValidatePath(dir string) bool {
	if dir == "" || pagesAll.ContainsPath(dir) {
		return false
	}
	return true
}

func getCurrentDirWithoutFile(dir string) string {
	if strings.HasSuffix(dir, "/") {
		dir = dir[:strings.LastIndex(dir, "/")]
	}
	if dir == "" {
		return "/"
	}
	lastIndexOfSlash := strings.LastIndex(dir, "/")
	if lastIndexOfSlash == len(dir)-1 {
		//Получили shash
		return dir
	}
	if strings.Contains(dir[lastIndexOfSlash:], ".") {
		return dir[:lastIndexOfSlash] + "/"
	}
	return dir + "/"
}

func concatenateUrlAndDir(url, dir string) string {
	if strings.HasSuffix(url, "/") && strings.HasPrefix(dir, "/") {
		// Doc: url = http://example.com/ && dir = /styles.css ||
		dir = dir[1:]
	} else if !strings.HasSuffix(url, "/") && !strings.HasPrefix(dir, "/") {
		// Doc: url = http://example.com && dir = styles.css
		dir = "/" + dir
	}
	return url + dir
}

func (ps pagesSlice) ContainsPath(path string) bool {
	for _, page := range ps {
		if page.path == path {
			return true
		}
	}
	return false
}

func (ps pagesSlice) DeleteDuplicates(anotherPagesSlice pagesSlice) pagesSlice {
	var pagesSliceToReturn pagesSlice
	for _, page := range ps {
		if anotherPagesSlice.ContainsPath(page.path) {
			pagesSliceToReturn = append(pagesSliceToReturn, page)
		}
	}
	return pagesSliceToReturn
}

func getAllSrcInCSSFile(contentCSS string) []string {
	var urls []string
	s := scanner.New(contentCSS)
	token := s.Next()
	for token.Type != scanner.TokenEOF {
		if token.Type == scanner.TokenURI {
			fmt.Println("Src: ", token.Value)
			url := strings.Trim(token.Value, "url(")
			url = strings.Trim(url, ")")
			urls = append(urls, url)
		}
		token = s.Next()
	}
	return urls
}
