package server

import (
	"net/http"
	"os"
	"recruitme/internal/routes"

	"github.com/a-h/templ"
	"github.com/gorilla/csrf"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	// Register routes
	fileServer := http.FileServer(http.FS(Files))
	mux.Handle("/assets/", fileServer)
	mux.Handle("/", templ.Handler(routes.Home()))
	mux.Handle("/scout", routes.NewHandler())
	mux.HandleFunc("/scoutme", routes.SendScoutMail)
	mux.Handle("/scouted", templ.Handler(routes.Scouted()))

	CSRF := csrf.Protect([]byte(os.Getenv("CSRF_KEY")))
	// Wrap the mux with CORS middleware and the CSRF middleware
	return CSRF(s.corsMiddleware(mux))
}

func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Replace "*" with specific origins if needed
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Set to "true" if credentials are required

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}
