// Package blog
package blog

import (
	"log/slog"
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
		templ.Handle(BlogPost()).ServeHTTP(w, r)
		// m := ContactModel{}
		// templ.Handler(Contact(&m, z.ZogIssueMap{})).ServeHTTP(w, r)

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func validateSchema(w http.ResponseWriter, r *http.Request) *ContactModel {

	model := ContactModel{
		Name:        r.PostForm.Get("name"),
		Email:       r.PostForm.Get("email"),
		Position:    r.PostForm.Get("position"),
		Level:       r.PostForm.Get("level"),
		Company:     r.PostForm.Get("company"),
		Description: r.PostForm.Get("description"),
		Link:        r.PostForm.Get("link"),
	}
	slog.Debug("Validating form")
	errs := ContactSchema.Validate(&model)
	// HTMX optimization using fragments.
	// Whenever there is an HX-Trigger-Name for this endpoint,
	// we know that there's a fragment of the same name and can send only that.
	if fragmentName := r.Header.Get("HX-Trigger-Name"); fragmentName != "" {
		slog.Info("Returning with fragments", "HX-Trigger-Name", fragmentName)
		templ.Handler(Contact(&model, errs), templ.WithFragments(fragmentName)).ServeHTTP(w, r)
		return nil
	}

	if len(errs) > 0 {
		slog.Info("Model was not parsed correctly: ", "Model Values", model, "Schema Errors", errs)
		templ.Handler(Contact(&model, errs)).ServeHTTP(w, r)
		return nil
	}
	return &model
}
