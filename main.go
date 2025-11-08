package main

import (
	"io"
	// "os"
	// "bytes"
	"fmt"
	// "log"
	"net/http"
	"html/template"
	// "github.com/AufaIzza/ASBC/api"
	"github.com/AufaIzza/ASBC/helper"
)

// func rootHandler(w http.ResponseWriter, req *http.Request) {
// 	if req.Method == http.MethodGet {
// 		io.WriteString(w, "Hello World\n")
// 	}
// }

var err error 

type Page struct {
	page string
}

func (p Page) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		io.WriteString(w, p.page)
	}
}

type Site struct {
	Mux *http.ServeMux
	Tmpl *template.Template
	Port string
}

func MakeSite(port string, partialsPath string) (*Site, error) {
	site := Site {
		Port: port,
	}
	site.Tmpl, err = helper.MakeTemplateFromGlobAndCollect("./site_src/partials/*.gohtml")
	if err != nil {
		return nil, err
	}
	site.Mux = http.NewServeMux()
	return &site, nil
}

func (s *Site) Route(path string, filepath string, props map[string]string) {
	err := helper.MakePageFromFile(s.Mux, s.Tmpl, path, filepath, props)
	if err != nil {
		panic(err)
	}
}

func (s *Site) Serve() error {
	fmt.Println("Listening on port ", s.Port)
	return http.ListenAndServe(s.Port, s.Mux)
}

func (s *Site) StaticPath(path string) {
	fs := http.FileServer(http.Dir(path))
	s.Mux.Handle("/static/", http.StripPrefix("/static/", fs))
}

func main() {
	// str := "<h1>this is not from partial</h1>{{template \"test\"}}"

	// site.Tmpl, err = helper.MakeTemplateFromGlobAndCollect("./site_src/partials/*.gohtml")
	// if err != nil {
	// 	panic(err)
	// }

	// tmpl, err := site.Tmpl.New("html").Parse(str)
	// if err != nil {
	// 	panic(err)
	// }

	// buf := bytes.NewBuffer(nil)
	// err = tmpl.Execute(buf, nil)
	// if err != nil {
	// 	panic(err)
	// }

	// page := Page {
	// 	page: buf.String(),
	// }

	// site.Mux.Handle("/", page)

	site, err := MakeSite(":6969", "./site_src/partials/*.gohtml")
	site.StaticPath("./site_src/static/")
	site.Route("/", "./site_src/pages/index.gohtml", map[string]string {
		"title": "ASBC - Main",
	})

	err = site.Serve()
	if err != nil {
		panic(err)
	}
}
