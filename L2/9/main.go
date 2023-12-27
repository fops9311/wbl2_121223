// © github.com/MikhailSolovev
package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type WebPage struct {
	uraw   string
	host   string
	scheme string
	doc    *goquery.Document
	// Множество всех ссылок
	links map[string]struct{}
}

// Блок assets, работает только при наличии --page-requisites

// FindAndDownloadImgs - позволяет найти все url картинок на странице (webpage), а затем скачать каждую из них
func (w *WebPage) FindAndDownloadImgs() {
	w.doc.Find("img").Each(func(index int, element *goquery.Selection) {
		imgSrc, exists := element.Attr("src")
		if exists {
			wrapDownloadFile(w.scheme + "://" + w.host + "/" + imgSrc)
		}
	})
}

// FindAndDownloadJS - позволяет найти все JS скрипты на странице (webpage), а затем скачать каждый из них. Скачивает
// только те, у которых тот же хост, что и у изначального сайта
func (w *WebPage) FindAndDownloadJS() {
	w.doc.Find("script").Each(func(index int, element *goquery.Selection) {
		scriptSrc, exists := element.Attr("src")
		if exists {
			if checkTwoRawUrls(w.uraw, scriptSrc) {
				wrapDownloadFile(scriptSrc)
			}
		}
	})
}

// FindAndDownloadCSS - позволяет найти все CSS на странице (webpage), а затем скачать каждый из них
func (w *WebPage) FindAndDownloadCSS() {
	w.doc.Find("link").Each(func(index int, element *goquery.Selection) {
		cssSrc, exists := element.Attr("href")
		if exists && cssSrc[len(cssSrc)-4:] == ".css" {
			wrapDownloadFile(w.scheme + "://" + w.host + "/" + cssSrc)
		}
	})
}

// FindAndDownloadICO - позволяет найти все ICO на странице (webpage), а затем скачать каждый из них
func (w *WebPage) FindAndDownloadICO() {
	w.doc.Find("link").Each(func(index int, element *goquery.Selection) {
		icoSrc, exists := element.Attr("href")
		if exists && icoSrc[len(icoSrc)-4:] == ".ico" {
			wrapDownloadFile(w.scheme + "://" + w.host + "/" + icoSrc)
		}
	})
}

// FindLinks - позволяет найти все url ссылок на странице (webpage), flag отвечает за --domains, переход только по тем
// ссылкам, которые относятся к нашему сайту
func (w *WebPage) FindLinks(flag bool) {
	w.doc.Find("a").Each(func(index int, element *goquery.Selection) {
		href, exists := element.Attr("href")
		if exists {
			if flag {
				// Сравниваем хосты страницы и ссылки, если они одинаковые, то добавляем в список ссылок
				if checkTwoRawUrls(w.uraw, href) {
					w.links[href] = struct{}{}
				} else if checkEmptyHost(href) {
					w.links[w.scheme+"://"+strings.Replace(w.host+"/"+href, "//", "/", -1)] = struct{}{}
				}
			} else {
				if checkEmptyHost(href) {
					w.links[w.scheme+"://"+strings.Replace(w.host+"/"+href, "//", "/", -1)] = struct{}{}
				} else {
					w.links[href] = struct{}{}
				}
			}
		}
	})
}

// checkEmptyHost - проверяет пустой ли хост у ссылки
func checkEmptyHost(uraw string) bool {
	u, err := url.Parse(uraw)
	if err != nil {
		log.Fatal(err)
	}

	return u.Host == ""
}

// checkTwoRawUrls - проверяет одинаковые ли хосты у двух ссылок
func checkTwoRawUrls(uraw1, uraw2 string) bool {
	u1, err := url.Parse(uraw1)
	if err != nil {
		log.Fatal(err)
	}

	u2, err := url.Parse(uraw2)
	if err != nil {
		log.Fatal(err)
	}

	return u1.Host == u2.Host
}

// DownloadHTML - позволяет скачать HTML страницы (web page), возвращает true, если страница уже скачана
func (w *WebPage) DownloadHTML() bool {
	return wrapDownloadFile(w.uraw)
}

// NewWebPage - конструктор для WebPage
func NewWebPage(uraw string) *WebPage {
	response := getResponse(uraw)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(response.Body)

	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	u, err := url.Parse(uraw)
	if err != nil {
		log.Fatal(err)
	}

	return &WebPage{
		uraw:   uraw,
		host:   u.Host,
		scheme: u.Scheme,
		doc:    doc,
		links:  map[string]struct{}{},
	}
}

// getResponse - GET запрос на заданный url, тело запроса не закрывается
func getResponse(uraw string) *http.Response {
	response, err := http.Get(uraw)
	if err != nil {
		log.Fatal(err)
	}

	if response.StatusCode > 400 {
		log.Fatal("Status code:", response.StatusCode)
	}

	return response
}

// downloadFile - скачивание любого файла, возвращает true, если файл существует
func downloadFile(uraw string, fileName string) bool {
	response := getResponse(uraw)
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(response.Body)

	// Создание каталогов до файла, если эти каталоги не существуют
	path := strings.Split(fileName, "/")
	err := os.MkdirAll(strings.Join(path[:len(path)-1], "/"), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	// Проверка существования файла
	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		return true
	}
	// Создание файла
	file, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(file)

	// Копирование тела запроса в файл
	_, err = io.Copy(file, response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return false
}

// wrapDownloadFile - обертка для downloadFile, которая сохраняет структуру относительно корня сайта
func wrapDownloadFile(uraw string) bool {
	u, err := url.Parse(uraw)
	if err != nil {
		log.Fatal(err)
	}
	if u.Path == "/" {
		return downloadFile(uraw, "./"+u.Host+u.Path+"index.html")
	}
	return downloadFile(uraw, "./"+u.Host+u.Path)
}

// Download - главная функция нашей программы
// uraw - строка, ссылка на webpage
// rec - флаг для рекурсивного скачивания
// assets - флаг для скачивания asset(ов)
// dom - флаг для скачивания только того, что принадлежит изначальному хосту
func Download(uraw string, rec, assets, dom bool) {
	page := NewWebPage(uraw)
	exists := page.DownloadHTML()
	if exists {
		return
	}
	// Скачивание assets страницы
	if assets {
		page.FindAndDownloadImgs()
		page.FindAndDownloadICO()
		page.FindAndDownloadCSS()
		page.FindAndDownloadJS()
	}
	if rec {
		page.FindLinks(dom)
		for link := range page.links {
			Download(link, rec, assets, dom)
		}
	}
}

// Use wget [-flags] url
func main() {
	// рекурсивоное скачивание (-r)
	r := flag.Bool("r", false, "recursive download webpage")
	// скачивание всех ассетов (-p)
	p := flag.Bool("p", false, "download ICO, CSS, JS, img")
	// проверка принадлежности ссылик изначальному хосту
	D := flag.Bool("D", false, "don’t follow any links outside of the website")

	flag.Parse()

	Download(flag.Arg(0), *r, *p, *D)
}
