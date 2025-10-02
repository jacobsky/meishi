package server

import (
	"crypto/rand"
	"log/slog"
	"net/http"
	"recruitme/internal/routes"
	"recruitme/internal/routes/contact"
	"strings"

	"github.com/a-h/templ"
	"github.com/gorilla/csrf"
	"github.com/invopop/ctxi18n"
)

func (s *Server) RegisterRoutes() http.Handler {
	// TODO: Switch the router to gorilla/mux for wildcard matching and subroutes
	// If loading fails, app should crash
	if err := ctxi18n.Load(Locales); err != nil {
		panic(err)
	}

	// Register routes
	webmux := http.NewServeMux()
	webmux.Handle("/", templ.Handler(routes.Home()))
	// TODO: Add micro blog functionality
	// mux.Handle("/blog", templ.Handler(routes.BlogPost()))
	webmux.Handle("/contact", contact.NewHandler())
	webmux.Handle("GET /contactcomplete", contact.NewHandler())

	rootmux := http.NewServeMux()
	fileServer := http.FileServer(http.FS(Files))

	rootmux.Handle("/assets/", fileServer)
	rootmux.Handle("/", http.RedirectHandler("/en/", http.StatusPermanentRedirect))
	rootmux.Handle("/en/", http.StripPrefix("/en", webmux))
	rootmux.Handle("/jp/", http.StripPrefix("/jp", webmux))

	// Handle middlewares
	csrfKey := generateCSRFKey()
	csrfMiddleware := csrf.Protect(
		csrfKey,
		csrf.TrustedOrigins(
			[]string{"localhost:8080", "127.0.0.1:8080"}),
		csrf.FieldName("_csrf"),
	)
	// Wrap the mux with CORS middleware and the CSRF middleware
	return csrfMiddleware(s.corsMiddleware(i18nmiddleware(rootmux)))
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
		// The local path will always be the first after "/" such as "/en-US/"
		slog.Debug("i18nmiddleware checking path", "path", r.URL.Path)
		pathSegments := strings.Split(r.URL.Path, "/")
		if len(pathSegments) > 1 {
			lang = pathSegments[1]
		} else {
			r.URL.Path = "/en" + r.URL.Path
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
