package server

import (
	"crypto/rand"
	"log/slog"
	"net/http"
	"recruitme/internal/routes"
	contact "recruitme/internal/routes/contactme"
	"strings"

	"github.com/a-h/templ"
	"github.com/gorilla/csrf"
	"github.com/invopop/ctxi18n"
)

func (s *Server) RegisterRoutes() http.Handler {
	// If loading fails, app should crash
	if err := ctxi18n.Load(Locales); err != nil {
		panic(err)
	}

	// Register routes
	mux := http.NewServeMux()
	fileServer := http.FileServer(http.FS(Files))
	mux.Handle("/assets/", fileServer)
	mux.Handle("/", templ.Handler(routes.Home()))
	// TODO: Add micro blog functionality
	// mux.Handle("/blog", templ.Handler(routes.BlogPost()))
	mux.Handle("/scout", contact.NewHandler())
	mux.HandleFunc("POST /scoutme", contact.SendScoutMail)
	mux.Handle("GET /scouted", templ.Handler(contact.Scouted()))
	middle := i18nmiddleware(mux)
	csrfKey := generateCSRFKey()
	csrfMiddleware := csrf.Protect(
		csrfKey,
		csrf.TrustedOrigins(
			[]string{"localhost:8080", "127.0.0.1:8080"}),
		csrf.FieldName("_csrf"),
	)
	// Wrap the mux with CORS middleware and the CSRF middleware
	return csrfMiddleware(s.corsMiddleware(middle))
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

func i18nmiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lang := "en" // default language
		pathSegments := strings.Split(r.URL.Path, "/")
		// The local path will always be the first after "/" such as "/en-US/"
		if len(pathSegments) > 1 {
			lang = pathSegments[1]
		}
		ctx, err := ctxi18n.WithLocale(r.Context(), lang)
		if err != nil {
			slog.Error("Error setting local", "Error", err)
			http.Error(w, "Error setting local", http.StatusBadRequest)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
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
