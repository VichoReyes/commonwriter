package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	content := r.URL.Query().Get("content")
	author := r.URL.Query().Get("author")
	title := r.URL.Query().Get("title")
	latestStory().Append(content, author, title)
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
	path := r.URL.Path
	story, err := lookupStory(path)
	if err != nil {
		http.Error(w, "Error 404: Story not found", http.StatusNotFound)
		return
	}
	templ.Execute(w, story)
}

func lookupStory(urlpath string) (*threads.Node, error) {
	// TODO join the two for loops, because creating a slice is redundant
	split := strings.Split(urlpath, "/")
	indices := make([]int, 0, len(split))
	for _, s := range split {
		if s == "" {
			continue
		}
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("lookupStory: %s can't be converted to int", s)
		}
		indices = append(indices, i)
	}

	n := &initialStory
	for _, i := range indices {
		n = n.Children()[i] // TODO out-of-bounds checking
		/*
			if !ok {
				return nil, fmt.Errorf("%v not a valid story path (yet)", path)
			}
		*/
	}
	return n, nil
}

func latestStory() *threads.Node {
	n := &initialStory
	for ; len(n.Children()) != 0; n = n.Children()[0] {
	}
	return n
}
