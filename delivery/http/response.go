package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"myapp/domain"
)

type errorBody struct {
	Error errorDetail `json:"error"`
}

type errorDetail struct {
	Code    string       `json:"code"`
	Message string       `json:"message"`
	Details []fieldError `json:"details,omitempty"`
}

type fieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func respondOK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

func respondCreated(c *gin.Context, data any) {
	c.JSON(http.StatusCreated, data)
}

func respondNoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func respondValidationError(c *gin.Context, field, msg string) {
	c.JSON(http.StatusBadRequest, errorBody{Error: errorDetail{
		Code:    "VALIDATION_ERROR",
		Message: "Input tidak valid",
		Details: []fieldError{{Field: field, Message: msg}},
	}})
}

func respondError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrEventNotFound):
		c.JSON(http.StatusNotFound, errorBody{Error: errorDetail{
			Code: "EVENT_NOT_FOUND", Message: "Event tidak ditemukan.",
		}})
	case errors.Is(err, domain.ErrRegistrationNotFound):
		c.JSON(http.StatusNotFound, errorBody{Error: errorDetail{
			Code: "REGISTRATION_NOT_FOUND", Message: "Registrasi tidak ditemukan.",
		}})
	case errors.Is(err, domain.ErrAdminNotFound):
		c.JSON(http.StatusNotFound, errorBody{Error: errorDetail{
			Code: "ADMIN_NOT_FOUND", Message: "Admin tidak ditemukan.",
		}})
	case errors.Is(err, domain.ErrDuplicateRegistration):
		c.JSON(http.StatusConflict, errorBody{Error: errorDetail{
			Code: "DUPLICATE_REGISTRATION", Message: "Email sudah terdaftar untuk event ini.",
		}})
	case errors.Is(err, domain.ErrInvalidStatusTransition):
		c.JSON(http.StatusBadRequest, errorBody{Error: errorDetail{
			Code: "INVALID_STATUS_TRANSITION", Message: "Transisi status tidak valid.",
		}})
	case errors.Is(err, domain.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, errorBody{Error: errorDetail{
			Code: "INVALID_CREDENTIALS", Message: "Email atau password salah.",
		}})
	case errors.Is(err, domain.ErrUnauthorized):
		c.JSON(http.StatusUnauthorized, errorBody{Error: errorDetail{
			Code: "UNAUTHORIZED", Message: "Token tidak valid atau sudah habis masa berlakunya.",
		}})
	case errors.Is(err, domain.ErrFileTooLarge):
		c.JSON(http.StatusRequestEntityTooLarge, errorBody{Error: errorDetail{
			Code: "FILE_TOO_LARGE", Message: "File bukti pembayaran terlalu besar (maks. 5 MB).",
		}})
	case errors.Is(err, domain.ErrInvalidMIMEType):
		c.JSON(http.StatusBadRequest, errorBody{Error: errorDetail{
			Code: "VALIDATION_ERROR", Message: "File harus berupa gambar (jpeg/png/webp).",
			Details: []fieldError{{Field: "payment_proof", Message: "File harus berupa gambar (jpeg/png/webp) maksimal 5 MB"}},
		}})
	default:
		c.JSON(http.StatusInternalServerError, errorBody{Error: errorDetail{
			Code: "INTERNAL_ERROR", Message: "Terjadi kesalahan pada server.",
		}})
	}
}
