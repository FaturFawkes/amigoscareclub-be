package http

import (
	"encoding/json"
	"net/http"

	"myapp/application/dto"
	"myapp/application/usecase"
)

// RegistrationHandler exposes HTTP endpoints for registrations.
type RegistrationHandler struct {
	placeUseCase *usecase.PlaceRegistrationUseCase
	getUseCase   *usecase.GetRegistrationUseCase
}

// NewRegistrationHandler creates a handler instance.
func NewRegistrationHandler(
	placeUseCase *usecase.PlaceRegistrationUseCase,
	getUseCase *usecase.GetRegistrationUseCase,
) *RegistrationHandler {
	return &RegistrationHandler{
		placeUseCase: placeUseCase,
		getUseCase:   getUseCase,
	}
}

// Place handles POST /registrations.
func (h *RegistrationHandler) Place(w http.ResponseWriter, r *http.Request) {
	var input dto.PlaceRegistrationInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	output, err := h.placeUseCase.Execute(r.Context(), input)
	if err != nil {
		http.Error(w, "failed to register", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}

// Get handles GET /registrations/{id}.
func (h *RegistrationHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "missing id", http.StatusBadRequest)
		return
	}

	output, err := h.getUseCase.Execute(r.Context(), id)
	if err != nil {
		http.Error(w, "registration not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(output)
}
