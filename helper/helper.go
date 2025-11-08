package helper

import (
	"fmt"
	"os"
	"html/template"
	"strings"
	"path/filepath"
	"net/http"
	"bytes"
	"io"
)

func basename(path string) string {
	file := filepath.Base(path)
	return strings.TrimSuffix(file, filepath.Ext(file))
}

func MakeTemplateFromFile(path string) (*template.Template, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	tmpl, err := template.New(basename(path)).Parse(string(data))
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func AddToTemplateFromFile(tmpl *template.Template, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	tmpl, err = tmpl.New(basename(path)).Parse(string(data))
	if err != nil {
		return err
	}
	return nil
}

func MakeTemplateFromGlobAndCollect(pattern string) (*template.Template, error) {
	paths, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	tmpl, err := MakeTemplateFromFile(paths[0])
	if err != nil { 
		return nil, err
	}	
	for idx, element := range paths {
		if idx != 0 {
			path := fmt.Sprintf("./%s", element)
			err = AddToTemplateFromFile(tmpl, path)
			if err != nil { 
				return nil, err
			}	
		}
	}
	return tmpl, err
}

func MakePageFromFile(
	mux *http.ServeMux, 
	base *template.Template, 
	path string, 
	filepath string,
	props map[string]string,
) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return err
	}
	file := basename(filepath)
	tmpl, err := base.New(file).Parse(string(data[:]))
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	err = tmpl.Execute(buf, props)
	if err != nil {
		return err
	}

	handler := func(w http.ResponseWriter, req *http.Request) {
		if req.Method == http.MethodGet {
			io.WriteString(w, buf.String())
		}
	}
	mux.HandleFunc(path, handler)
	return nil
}
