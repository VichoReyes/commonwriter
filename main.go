package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
	"syscall"

	"github.com/esclerofilo/commonwriter/threads"
)

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/upload/", upload)
	http.HandleFunc("/images/", serveImage)
	// log.Fatal(http.ListenAndServe("0.0.0.0:3000", nil)) // "prod"
	log.Fatal(http.ListenAndServe("localhost:8000", nil)) // dev
}

func upload(w http.ResponseWriter, r *http.Request) {
	content := r.FormValue("content")
	author := r.FormValue("author")
	title := r.FormValue("title")
	if content == "" || author == "" || title == "" {
		w.Write([]byte("one or more fields was empty"))
		return
	}

	// skip the "/upload/" part
	wholePath := strings.Split(r.URL.Path, "/")
	storyCoords := strings.Join(wholePath[2:], "/")

	story, err := lookupStory(storyCoords)
	if err != nil {
		http.Error(w, "Error 404: Story not found", http.StatusNotFound)
		return
	}
	newIndex := story.Append(content, author, title)
	// TODO change this to a redirect
	context := struct {
		NewURL string
		Title  string
	}{
		path.Clean("/" + storyCoords + "/" + strconv.Itoa(newIndex)),
		title,
	}
	templ.ExecuteTemplate(w, "uploadSuccessful.html", context)
}

var initialStory threads.Node

var templ = template.Must(template.ParseFiles("base.html", "edit.html", "uploadSuccessful.html"))

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
	context := struct {
		*http.Request
		*threads.Node
	}{
		r,
		story,
	}
	switch r.URL.Query().Get("view") {
	//	case "append":
	//		templ.ExecuteTemplate(w, "append.html", story)
	case "edit":
		templ.ExecuteTemplate(w, "edit.html", context)
	default:
		templ.ExecuteTemplate(w, "base.html", context)
	}
}

func lookupStory(urlpath string) (*threads.Node, error) {
	split := strings.Split(urlpath, "/")
	n := &initialStory
	for _, s := range split {
		if s == "" {
			continue
		}
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("lookupStory: %s can't be converted to int", s)
		}
		var ok bool
		n, ok = n.Child(i)
		if !ok {
			return nil, fmt.Errorf("lookupStory: index %d too large", i)
		}
	}
	return n, nil
}
