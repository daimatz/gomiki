package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/russross/blackfriday"
)

const Home string = "$HOME/Dropbox"
const Addr string = ":8080"

type Config struct {
	Home string
	Port int
}

type RootHandler struct{}

func (h RootHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var u *url.URL = req.URL
	var path string = os.ExpandEnv(Home) + u.Path

	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		log.Printf("failed to open %s\n", path)
	}

	md, err := ioutil.ReadAll(file)
	if err != nil {
		log.Printf("failed to read %s\n", path)
	}

	html := blackfriday.MarkdownCommon([]byte(md))
	w.Write([]byte(html))
}

func main() {
	server := &http.Server{
		Addr:    Addr,
		Handler: RootHandler{},
	}
	log.Fatal(server.ListenAndServe())
}
