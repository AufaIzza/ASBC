package api

import (
	"io"
	"net/http"
)

func TestHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		io.WriteString(w, "I AM Testing")
	}
}
