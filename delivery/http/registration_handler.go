package http

import (
	"mime/multipart"

	"github.com/gin-gonic/gin"
	"myapp/application/dto"
	"myapp/application/usecase"
)

// RegistrationHandler handles public registration endpoints.
type RegistrationHandler struct {
	createUC *usecase.CreateRegistrationUseCase
	getUC    *usecase.GetRegistrationUseCase
}

// NewRegistrationHandler wires the use cases.
func NewRegistrationHandler(
	createUC *usecase.CreateRegistrationUseCase,
	getUC *usecase.GetRegistrationUseCase,
) *RegistrationHandler {
	return &RegistrationHandler{createUC: createUC, getUC: getUC}
}

type createRegistrationRequest struct {
	Name  string `json:"name"  form:"name"`
	Email string `json:"email" form:"email"`
	Phone string `json:"phone" form:"phone"`
	Age   int    `json:"age"   form:"age"`
}

// Create handles POST /events/:eventSlug/registrations.
func (h *RegistrationHandler) Create(c *gin.Context) {
	var req createRegistrationRequest
	if err := c.ShouldBind(&req); err != nil {
		respondValidationError(c, "payload", "request tidak valid")
		return
	}

	type reqField struct{ name, val string }
	for _, f := range []reqField{
		{"name", req.Name},
		{"email", req.Email},
		{"phone", req.Phone},
	} {
		if f.val == "" {
			respondValidationError(c, f.name, f.name+" wajib diisi")
			return
		}
	}
	if req.Age < 10 || req.Age > 100 {
		respondValidationError(c, "age", "Usia harus berupa angka antara 10 dan 100")
		return
	}

	var paymentProof *multipart.FileHeader
	if fh, err := c.FormFile("payment_proof"); err == nil {
		paymentProof = fh
	}

	out, err := h.createUC.Execute(c.Request.Context(), dto.CreateRegistrationInput{
		EventSlug:    c.Param("eventSlug"),
		Name:         req.Name,
		Email:        req.Email,
		Phone:        req.Phone,
		Age:          req.Age,
		PaymentProof: paymentProof,
	})
	if err != nil {
		respondError(c, err)
		return
	}

	respondCreated(c, out)
}

// Get handles GET /events/:eventSlug/registrations/:registrationId.
func (h *RegistrationHandler) Get(c *gin.Context) {
	out, err := h.getUC.Execute(c.Request.Context(), dto.GetRegistrationInput{
		EventSlug:      c.Param("eventSlug"),
		RegistrationID: c.Param("registrationId"),
	})
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, out)
}
