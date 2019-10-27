package handlers

import (
	"github.com/dmoles/adler/server/util"
	"github.com/gobuffalo/packr"
	"github.com/gorilla/mux"
	"net/http"
)

type boxHandler struct {
	box packr.Box
	varname string
}

func (h *boxHandler) Handle(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	path := vars[h.varname]
	if path == "" {
		http.NotFound(w, r)
		return
	}
	data, err := h.box.Find(path)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	util.WriteData(w, path, data)
}


