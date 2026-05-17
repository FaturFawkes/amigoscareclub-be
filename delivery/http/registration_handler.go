package http

import (
	"strconv"

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

// Create handles POST /events/:eventSlug/registrations (multipart/form-data).
func (h *RegistrationHandler) Create(c *gin.Context) {
	if err := c.Request.ParseMultipartForm(6 << 20); err != nil {
		respondValidationError(c, "payload", "invalid multipart form")
		return
	}

	name := c.PostForm("name")
	email := c.PostForm("email")
	phone := c.PostForm("phone")
	ageStr := c.PostForm("age")
	coffeeChoice := c.PostForm("coffee_choice")

	type reqField struct{ name, val string }
	for _, f := range []reqField{
		{"name", name}, {"email", email}, {"phone", phone}, {"age", ageStr}, {"coffee_choice", coffeeChoice},
	} {
		if f.val == "" {
			respondValidationError(c, f.name, f.name+" wajib diisi")
			return
		}
	}

	age, err := strconv.Atoi(ageStr)
	if err != nil || age < 10 || age > 100 {
		respondValidationError(c, "age", "Usia harus berupa angka antara 10 dan 100")
		return
	}

	files := c.Request.MultipartForm.File["payment_proof"]
	if len(files) == 0 {
		respondValidationError(c, "payment_proof", "Bukti pembayaran wajib diunggah")
		return
	}

	out, err := h.createUC.Execute(c.Request.Context(), dto.CreateRegistrationInput{
		EventSlug:    c.Param("eventSlug"),
		Name:         name,
		Email:        email,
		Phone:        phone,
		Age:          age,
		CoffeeChoice: coffeeChoice,
		PaymentProof: files[0],
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
