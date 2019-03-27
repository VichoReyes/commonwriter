package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"syscall"

	"github.com/esclerofilo/commonwriter/threads"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/story/{key}", StoryHandler)
	r.HandleFunc("/upload/{key}", UploadHandler)
	r.HandleFunc("/images/", serveImage)
	// log.Fatal(http.ListenAndServe("0.0.0.0:3000", r)) // "prod"
	log.Fatal(http.ListenAndServe("localhost:8080", r)) // dev
}

// HomeHandler is a handler for GET /
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("This is the home."))
}

// UploadHandler is a handler for POST /upload/{key}
func UploadHandler(w http.ResponseWriter, r *http.Request) {
	content := r.FormValue("content")
	title := r.FormValue("title")
	if content == "" || title == "" {
		w.Write([]byte("one or more fields was empty"))
		return
	}

	keystring, ok := mux.Vars(r)["key"]
	if !ok {
		log.Panic("key missing")
	}
	key, err := strconv.ParseInt(keystring, 10, 64)
	if err != nil {
		log.Panic(err)
	}

	story, err := threads.Get(key)
	if err != nil {
		http.Error(w, "Error 404: Story not found", http.StatusNotFound)
		return
	}
	newID := story.Append(content, "", title)
	// TODO change this to a redirect
	context := struct {
		NewURL string
		Title  string
	}{
		"/story/" + strconv.FormatInt(newID, 10),
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

// StoryHandler is a handler for GET /story/{key}
func StoryHandler(w http.ResponseWriter, r *http.Request) {
	keystring, ok := mux.Vars(r)["key"]
	if !ok {
		log.Panic("key missing")
	}
	key, err := strconv.ParseInt(keystring, 10, 64)
	if err != nil {
		log.Panic(err)
	}

	story, err := threads.Get(key)
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
