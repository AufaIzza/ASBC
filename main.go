package main

import (
	"io"
	// "os"
	"bytes"
	// "fmt"
	"log"
	"net/http"
	// "html/template"
	// "github.com/AufaIzza/asbc/api"
	"github.com/AufaIzza/asbc/site"
)

// func rootHandler(w http.ResponseWriter, req *http.Request) {
// 	if req.Method == http.MethodGet {
// 		io.WriteString(w, "Hello World\n")
// 	}
// }

type Page struct {
	page string
}

func (p Page) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		io.WriteString(w, p.page)
	}
}

func main() {
	str := "<h1>this is not from partial</h1>{{template \"test\"}}"
	tmpl, err := site.MakeTemplatesFromGlobAndCollect("./site_src/partials/*.gohtml")
	if err != nil {
		panic(err)
	}
	tmpl, err = tmpl.New("html").Parse(str)
	if err != nil {
		panic(err)
	}

	// buf := new(bytes.Buffer)
	buf := bytes.NewBuffer(nil)
	err = tmpl.Execute(buf, nil)
	if err != nil {
		panic(err)
	}	
	page := Page {
		page: buf.String(),
	}
	http.Handle("/", page)
	log.Fatal(http.ListenAndServe(":6969", nil))	
}
