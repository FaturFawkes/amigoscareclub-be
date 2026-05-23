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
	listUC          *usecase.AdminListRegistrationsUseCase
	verifyUC        *usecase.AdminVerifyRegistrationUseCase
	resendTicketUC  *usecase.AdminResendTicketUseCase
	resendAllUC     *usecase.AdminResendAllTicketsUseCase
}

// NewAdminRegistrationHandler wires the use cases.
func NewAdminRegistrationHandler(
	listUC *usecase.AdminListRegistrationsUseCase,
	verifyUC *usecase.AdminVerifyRegistrationUseCase,
	resendTicketUC *usecase.AdminResendTicketUseCase,
	resendAllUC *usecase.AdminResendAllTicketsUseCase,
) *AdminRegistrationHandler {
	return &AdminRegistrationHandler{
		listUC:         listUC,
		verifyUC:       verifyUC,
		resendTicketUC: resendTicketUC,
		resendAllUC:    resendAllUC,
	}
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

// ResendTicket handles POST /admin/events/:eventSlug/registrations/:registrationId/resend-ticket.
func (h *AdminRegistrationHandler) ResendTicket(c *gin.Context) {
	out, err := h.resendTicketUC.Execute(c.Request.Context(), dto.ResendTicketInput{
		EventSlug:      c.Param("eventSlug"),
		RegistrationID: c.Param("registrationId"),
	})
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, out)
}

// ResendAllTickets handles POST /admin/events/:eventSlug/registrations/resend-all-tickets.
func (h *AdminRegistrationHandler) ResendAllTickets(c *gin.Context) {
	out, err := h.resendAllUC.Execute(c.Request.Context(), dto.ResendAllTicketsInput{
		EventSlug: c.Param("eventSlug"),
	})
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, out)
}
