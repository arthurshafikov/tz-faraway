package handler

import (
	"io"
	"net/http"
)

func NewHandler() *http.ServeMux {
	handler := &http.ServeMux{}

	handler.Handle("/", http.HandlerFunc(handleOK))

	return handler
}

func handleOK(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "OK\n")
}
