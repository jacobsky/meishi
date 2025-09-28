package routes

import (
	"net/http"

	z "github.com/Oudwins/zog"
	"github.com/a-h/templ"
)

type ScoutMeModel struct {
	Name        string
	Email       string
	Position    string
	Level       string
	CompanyName string
	Description string
	Link        string
}

var ScoutSchema = z.Struct(z.Shape{
	"name":         z.String().Required().Min(4),
	"email":        z.String().Required().Email(),
	"position":     z.String().Required(),
	"level":        z.String().Required(),
	"company_name": z.String().Required(),
	"description":  z.String().Required().Min(10),
	"link":         z.String().Required().URL(),
},
)

type Handler struct{}

func NewHandler() http.Handler {
	return &Handler{}
}
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		m := ScoutMeModel{}
		v := templ.Handler(Scout(&m, z.ZogIssueMap{}))
		v.ServeHTTP(w, r)
	case http.MethodPost:

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
