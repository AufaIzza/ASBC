package api

import (
	"io"
	"net/http"
	"strings"
	"fmt"
)

func AllPublicNotesHandlerTest(w http.ResponseWriter, r *http.Request) {
	var sb strings.Builder
	
	notes, err := queryAllPublicNotes()
	if err != nil {
		panic(err)
	}

	for _, note := range notes {
		sb.WriteString(fmt.Sprintf("%v", note))
	}

	io.WriteString(w, sb.String())
}

func TestHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		io.WriteString(w, "I AM Testing")
	}
}
