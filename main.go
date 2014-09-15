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
	var err error

	var u *url.URL = req.URL
	var localPath string = os.ExpandEnv(Home) + u.Path

	var file *os.File
	file, err = os.Open(localPath)
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

	var stat os.FileInfo
	stat, err = file.Stat()
	if err != nil {
		log.Printf("failed to get info %s\n", localPath)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	if stat.IsDir() {
		h.ServeDirectory(u.Path, localPath, w)
	} else {
		var ext string
		ext = path.Ext(localPath)
		if ext == ".md" || ext == ".markdown" || ext == ".mkd" || ext == ".mdown" {
			h.ServeMarkdown(file, u.Path, localPath, w)
		} else {
			var bytes []byte
			bytes, err = ioutil.ReadFile(localPath)
			var mimeType string = http.DetectContentType(bytes)
			w.Header().Set("Content-Type", mimeType)
			var wrote int
			wrote, err = w.Write(bytes)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			if wrote != len(bytes) {
				http.Error(w, "wrote bytes is not equal to file length", http.StatusInternalServerError)
			}
		}
	}
}

func (h *RootHandler) ServeMarkdown(file *os.File, requestPath string, localPath string, w http.ResponseWriter) {
	var err error

	var md []byte
	md, err = ioutil.ReadAll(file)
	if err != nil {
		log.Printf("failed to read %s\n", localPath)
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}
	h.WriteResponse(md, w)
}

func (h *RootHandler) ServeDirectory(requestPath string, localPath string, w http.ResponseWriter) {
	var err error

	var infos []os.FileInfo
	infos, err = ioutil.ReadDir(localPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	var buf bytes.Buffer
	var info os.FileInfo
	for _, info = range infos {
		var renderPath string
		renderPath = path.Join(requestPath, info.Name())
		buf.WriteString(fmt.Sprintf("- [%s](%s)\n", info.Name(), renderPath))
	}

	h.WriteResponse(buf.Bytes(), w)
}

func (h *RootHandler) WriteResponse(content []byte, w http.ResponseWriter) {
	var err error

	var str = html.MainArg{
		html.MarkdownToTemplateHTML(content),
	}
	err = html.MainTemplate.Execute(w, str)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func main() {
	var server *http.Server = &http.Server{
		Addr:    Addr,
		Handler: &RootHandler{},
	}
	log.Fatal(server.ListenAndServe())
}
