package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
)

// ------------------------------
// Exported symbols

type Handler interface {
	Register(r *mux.Router)
}

// ------------------------------
// Unexported symbols

type handler interface {
	Handler
	handle(w http.ResponseWriter, r *http.Request)
}