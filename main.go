package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"syscall"

	"github.com/esclerofilo/commonwriter/threads"
)

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/images/", serveImage)
	// log.Fatal(http.ListenAndServe("0.0.0.0:3000", nil)) // "prod"
	log.Fatal(http.ListenAndServe("localhost:8000", nil)) // dev
}

func upload(w http.ResponseWriter, r *http.Request) {
	latestStory().Append(r.URL.Query().Get("add"))
	log.Printf("%s", latestStory())
	w.Write([]byte("uploaded succesfully"))
}

var initialStory threads.Node

var templ = template.Must(template.ParseFiles("base.html"))

func serveImage(w http.ResponseWriter, r *http.Request) {
	// TODO url validation
	p := omitSlash(r.URL.Path)
	img, err := os.Open(p)
	if err != nil {
		err := err.(*os.PathError)
		if err.Err == syscall.ENOENT {
			http.Error(w, "image not found", http.StatusNotFound)
			log.Printf("failed attempt to retrieve image in %s", p)
			return
		}
		panic(err)
	}
	defer img.Close()
	io.Copy(w, img)
}

func omitSlash(original string) (fixed string) {
	if original[0] == '/' {
		return original[1:]
	}
	return original
}

func home(w http.ResponseWriter, r *http.Request) {
	templ.Execute(w, latestStory().String())
}

func latestStory() *threads.Node {
	n := &initialStory
	for ; len(n.Children()) != 0; n = n.Children()[0] {
	}
	return n
}
