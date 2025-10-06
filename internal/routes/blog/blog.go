// Package blog
package blog

import (
	"net/http"

	"github.com/a-h/templ"
)

type Handler struct{}

func NewHandler() http.Handler {
	return &Handler{}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		templ.Handler(BlogPost()).ServeHTTP(w, r)
		// m := ContactModel{}
		// templ.Handler(Contact(&m, z.ZogIssueMap{})).ServeHTTP(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
