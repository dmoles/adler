package handlers

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/dmoles/adler/resources"
)

type resourceHandler struct {
}

func ResourceHandler() Handler {
	return &resourceHandler{}
}

func (h *resourceHandler) Register(r *mux.Router) {
	r.MatcherFunc(h.isResource).HandlerFunc(h.handle)
}

func (h *resourceHandler) isResource(r *http.Request, _ *mux.RouteMatch) bool {
	urlPath := r.URL.Path
	_, err := resources.Resolve(urlPath)
	return err == nil
}

func (h *resourceHandler) handle(w http.ResponseWriter, r *http.Request) {
	if err := writeResource(w, r); err != nil {
		http.NotFound(w, r)
		return
	}
}

func writeResource(w http.ResponseWriter, r *http.Request) error {
	urlPath := r.URL.Path
	// log.Printf("writeResource(): %v", urlPath)

	resource, err := resources.Resolve(urlPath)
	if err != nil {
		return err
	}

	// If this fails, we've already started writing the response, so it's too
	// late to return a 404 or whatever; just log it and move on
	err = resource.Write(w, urlPath)
	if err != nil {
		log.Print(err)
	}
	return nil
}
