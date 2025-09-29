package server

import (
	"crypto/rand"
	"net/http"
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
	mux.HandleFunc("POST /scoutme", routes.SendScoutMail)
	mux.Handle("GET /scouted", templ.Handler(routes.Scouted()))
	csrfKey := generateCSRFKey()
	csrfMiddleware := csrf.Protect(
		csrfKey,
		csrf.TrustedOrigins(
			[]string{"localhost:8080", "127.0.0.1:8080"}),
		csrf.FieldName("_csrf"),
	)
	// Wrap the mux with CORS middleware and the CSRF middleware
	return csrfMiddleware(s.corsMiddleware(mux))
}
func generateCSRFKey() []byte {
	// Make a random 32 bit key for this application run.
	key := make([]byte, 32)
	n, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	if n != 32 {
		panic("Could not properly generate CSRF key")
	}
	return key
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
