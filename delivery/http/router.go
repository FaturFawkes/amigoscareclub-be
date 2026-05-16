package http

import "net/http"

// NewRouter wires HTTP routes to handlers.
func NewRouter(handler *RegistrationHandler) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/registrations", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			handler.Place(w, r)
		case http.MethodGet:
			handler.Get(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})

	return mux
}
