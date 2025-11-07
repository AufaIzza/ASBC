package site

import (
	"fmt"
	"os"
	"html/template"
	"strings"
	"path/filepath"
)

func basename(filename string) string {
	file := filepath.Base(filename)
	return strings.TrimSuffix(file, filepath.Ext(file))
}

func MakeTemplateFromFile(filename string) (*template.Template, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	tmpl, err := template.New(basename(filename)).Parse(string(data))
	if err != nil {
		return nil, err
	}
	return tmpl, nil
}

func MakeTemplatesFromGlobAndCollect(pattern string) (*template.Template, error) {
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
			tmpl, err = MakeTemplateFromFile(path)
			if err != nil { 
				return nil, err
			}	
		}
	}
	return tmpl, err
}
