package handler

import (
	"fmt"
	"net/http"
)

type handler struct{}

func New() *handler {
	return &handler{}
}

func (h *handler) Landing(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, Bean!")
}
