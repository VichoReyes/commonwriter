package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"syscall"

	"gopl.io/ch7/eval"
)

func main() {
	http.HandleFunc("/index", index)
	http.HandleFunc("/calculate", calculate)
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

func index(w http.ResponseWriter, r *http.Request) {

	templ.Execute(w, nil)
}

func parseAndCheck(s string) (eval.Expr, error) {
	if s == "" {
		return nil, fmt.Errorf("empty expression")
	}
	expr, err := eval.Parse(s)
	if err != nil {
		return nil, err
	}
	vars := make(map[eval.Var]bool)
	if err := expr.Check(vars); err != nil {
		return nil, err
	}
	if len(vars) != 0 {
		return nil, fmt.Errorf("expression had variables")
	}
	return expr, nil
}

func calculate(w http.ResponseWriter, r *http.Request) {
	expr, err := parseAndCheck(r.FormValue("expr"))
	if err != nil {
		http.Error(w, "bad expr: "+err.Error(), http.StatusBadRequest)
		return
	}

	res := expr.Eval(nil)
	fmt.Fprintf(w, "%f", res)
}
