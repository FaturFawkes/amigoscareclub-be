package http

import (
	"github.com/gin-gonic/gin"
	"myapp/application/dto"
	"myapp/application/usecase"
)

// EventHandler handles public event endpoints.
type EventHandler struct {
	getUC *usecase.GetEventUseCase
}

// NewEventHandler wires the use case.
func NewEventHandler(getUC *usecase.GetEventUseCase) *EventHandler {
	return &EventHandler{getUC: getUC}
}

// Get handles GET /events/:eventSlug.
func (h *EventHandler) Get(c *gin.Context) {
	out, err := h.getUC.Execute(c.Request.Context(), dto.GetEventInput{Slug: c.Param("eventSlug")})
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, out)
}
