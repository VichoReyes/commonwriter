package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"syscall"
)

func main() {
	currentStory = new(string)
	http.HandleFunc("/", home)
	http.HandleFunc("/upload", upload)
	http.HandleFunc("/images/", serveImage)
	// log.Fatal(http.ListenAndServe("0.0.0.0:3000", nil)) // "prod"
	log.Fatal(http.ListenAndServe("localhost:8000", nil)) // dev
}

// TODO
func upload(w http.ResponseWriter, r *http.Request) {
	storyMutex.Lock()
	*currentStory += r.URL.Query().Get("add")
	storyMutex.Unlock()
	w.Write([]byte("uploaded succesfully"))
}

// TODO: automate this somehow
var currentStory *string
var storyMutex sync.Mutex

var templ = template.Must(template.New("calc").Parse(`<html>
<body>
	<p>
	{{.}}
	</p>
</body>
</html>`))

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
	templ.Execute(w, *currentStory)
}
