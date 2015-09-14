package main

import (
	"flag"
	"fmt"
	"github.com/prodis/shortener/url"
	"log"
	"net/http"
	"strings"
)

var (
	port    *int
	baseUrl string
)

func init() {
	domain := flag.String("d", "localhost", "domain")
	port = flag.Int("p", 8888, "port")

	flag.Parse()

	baseUrl = fmt.Sprintf("http://%s:%d", *domain, *port)
}

func main() {
	url.ConfigureRepository(url.NewMemoryRepository())

	http.HandleFunc("/api/shorten", Shortener)
	http.HandleFunc("/r/", Redirecter)

	logInfo("Starting server using %d port...", *port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}

type Headers map[string]string

func Shortener(writer http.ResponseWriter, request *http.Request) {
	if request.Method != "POST" {
		respondWith(writer, http.StatusMethodNotAllowed, Headers{"Allow": "POST"})
		return
	}

	url, newUrl, err := url.FindOrCreate(extractUrl(request))

	if err != nil {
		respondWith(writer, http.StatusBadRequest, nil)
		return
	}

	var status int
	if newUrl {
		status = http.StatusCreated
	} else {
		status = http.StatusOK
	}

	shortUrl := fmt.Sprintf("%s/r/%s", baseUrl, url.Id)

	respondWith(writer, status, Headers{"Location": shortUrl})
	logInfo("URL %s shortened to %s", url.Target, shortUrl)
}

func Redirecter(writer http.ResponseWriter, request *http.Request) {
	path := strings.Split(request.URL.Path, "/")
	id := path[len(path)-1]

	if url := url.Find(id); url != nil {
		http.Redirect(writer, request, url.Target, http.StatusMovedPermanently)
		logInfo("URL %s redirected to %s", request.URL.Path, url.Target)
	} else {
		http.NotFound(writer, request)
	}
}

func extractUrl(request *http.Request) string {
	url := make([]byte, request.ContentLength, request.ContentLength)
	request.Body.Read(url)
	return string(url)
}

func respondWith(writer http.ResponseWriter, status int, headers Headers) {
	for key, value := range headers {
		writer.Header().Set(key, value)
	}

	writer.WriteHeader(status)
}

func logInfo(format string, values ...interface{}) {
	log.Printf(fmt.Sprintf("%s\n", format), values...)
}
