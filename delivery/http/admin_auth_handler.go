package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"myapp/application/dto"
	"myapp/application/serviceInterface"
	"myapp/application/usecase"
	"myapp/delivery/middleware"
)

// AdminAuthHandler handles admin authentication endpoints.
type AdminAuthHandler struct {
	loginUC  *usecase.AdminLoginUseCase
	meUC     *usecase.AdminMeUseCase
	logoutUC *usecase.AdminLogoutUseCase
}

// NewAdminAuthHandler wires the use cases.
func NewAdminAuthHandler(
	loginUC *usecase.AdminLoginUseCase,
	meUC *usecase.AdminMeUseCase,
	logoutUC *usecase.AdminLogoutUseCase,
) *AdminAuthHandler {
	return &AdminAuthHandler{loginUC: loginUC, meUC: meUC, logoutUC: logoutUC}
}

// Login handles POST /admin/auth/login.
func (h *AdminAuthHandler) Login(c *gin.Context) {
	var input dto.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil || input.Email == "" || input.Password == "" {
		c.JSON(http.StatusUnprocessableEntity, errorBody{Error: errorDetail{
			Code: "VALIDATION_ERROR", Message: "email dan password wajib diisi",
		}})
		return
	}

	out, err := h.loginUC.Execute(c.Request.Context(), input)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, out)
}

// Me handles GET /admin/auth/me.
func (h *AdminAuthHandler) Me(c *gin.Context) {
	claims := c.MustGet(middleware.AdminClaimsKey).(serviceInterface.TokenClaims)

	out, err := h.meUC.Execute(c.Request.Context(), claims.Sub)
	if err != nil {
		respondError(c, err)
		return
	}
	respondOK(c, out)
}

// Logout handles POST /admin/auth/logout.
func (h *AdminAuthHandler) Logout(c *gin.Context) {
	raw := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	if err := h.logoutUC.Execute(c.Request.Context(), raw); err != nil {
		respondError(c, err)
		return
	}
	respondNoContent(c)
}
