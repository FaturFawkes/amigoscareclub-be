package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"myapp/application/dto"
	"myapp/application/usecase"
)

// AdminRegistrationHandler handles admin-only registration management endpoints.
type AdminRegistrationHandler struct {
	listUC   *usecase.AdminListRegistrationsUseCase
	verifyUC *usecase.AdminVerifyRegistrationUseCase
}

// NewAdminRegistrationHandler wires the use cases.
func NewAdminRegistrationHandler(
	listUC *usecase.AdminListRegistrationsUseCase,
	verifyUC *usecase.AdminVerifyRegistrationUseCase,
) *AdminRegistrationHandler {
	return &AdminRegistrationHandler{listUC: listUC, verifyUC: verifyUC}
}

// List handles GET /admin/events/:eventSlug/registrations.
func (h *AdminRegistrationHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "20"))

	out, err := h.listUC.Execute(c.Request.Context(), dto.ListRegistrationsInput{
		EventSlug: c.Param("eventSlug"),
		Status:    c.Query("status"),
		Page:      page,
		PerPage:   perPage,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, out)
}

// Verify handles PATCH /admin/events/:eventSlug/registrations/:registrationId/verify.
func (h *AdminRegistrationHandler) Verify(c *gin.Context) {
	var body struct {
		Status string `json:"status"`
		Note   string `json:"note"`
	}
	if err := c.ShouldBindJSON(&body); err != nil || body.Status == "" {
		c.JSON(http.StatusBadRequest, errorBody{Error: errorDetail{
			Code: "VALIDATION_ERROR", Message: "Field 'status' wajib diisi",
			Details: []fieldError{{Field: "status", Message: "status wajib diisi"}},
		}})
		return
	}

	out, err := h.verifyUC.Execute(c.Request.Context(), dto.VerifyRegistrationInput{
		EventSlug:      c.Param("eventSlug"),
		RegistrationID: c.Param("registrationId"),
		Status:         body.Status,
		Note:           body.Note,
	})
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, out)
}
