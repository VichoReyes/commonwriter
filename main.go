package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"syscall"
)

func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/images/", serveImage)
	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

var templ = template.Must(template.New("calc").Parse(`<html>
<body>
	<img src="/images/1.jpg">
</body>
</html>`))

func serveImage(w http.ResponseWriter, r *http.Request) {
	// TODO url validation
	url := r.URL
	log.Println("url is", url)
	p := omitSlash(url.Path)
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

	templ.Execute(w, nil)
}
