package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"

	"github.com/daimatz/gomiki/html"
)

const Home string = "$HOME/Dropbox"
const Addr string = ":8080"

type RootHandler struct{}

func (h *RootHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var u *url.URL = req.URL
	var localPath string = os.ExpandEnv(Home) + u.Path

	file, err := os.Open(localPath)
	defer file.Close()
	if err != nil {
		log.Printf("failed to open %s\n", localPath)
		if os.IsNotExist(err) {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	stat, err := file.Stat()
	if err != nil {
		log.Printf("failed to get info %s\n", localPath)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	if stat.IsDir() {
		h.ServeDirectory(localPath, w)
	} else {
		path.Ext(localPath)
		md, err := ioutil.ReadAll(file)
		if err != nil {
			log.Printf("failed to read %s\n", localPath)
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		h.WriteResponse(md, w)
	}
}

func (h *RootHandler) ServeDirectory(localPath string, w http.ResponseWriter) {
	var buf bytes.Buffer

	infos, err := ioutil.ReadDir(localPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	for _, info := range infos {
		buf.WriteString(fmt.Sprintf("- [%s](%s)\n", info.Name, info.Name))
	}

	h.WriteResponse(buf.Bytes(), w)
}

func (h *RootHandler) WriteResponse(content []byte, w http.ResponseWriter) {
	err := html.MainTemplate.Execute(w, html.MainArg{
		html.MarkdownToTemplateHTML(content),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	server := &http.Server{
		Addr:    Addr,
		Handler: &RootHandler{},
	}
	log.Fatal(server.ListenAndServe())
}
