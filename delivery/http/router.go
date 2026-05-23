package http

import (
	"myapp/application/serviceInterface"
	"myapp/delivery/middleware"
	"myapp/domain"

	"github.com/gin-gonic/gin"
)

// NewRouter builds and returns the Gin engine with all routes registered under basePath.
func NewRouter(
	basePath string,
	tokenSvc serviceInterface.TokenService,
	tokenRepo domain.TokenRepository,
	regHandler *RegistrationHandler,
	eventHandler *EventHandler,
	adminAuthHandler *AdminAuthHandler,
	adminRegHandler *AdminRegistrationHandler,
) *gin.Engine {
	r := gin.New()
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())
	r.Static("/public", "./public")

	v1 := r.Group(basePath)

	// Public endpoints
	v1.GET("/events/:eventSlug", eventHandler.Get)
	v1.GET("/events/:eventSlug/registrations/:registrationId", regHandler.Get)
	v1.POST("/events/:eventSlug/registrations", regHandler.Create)
	v1.POST("/admin/auth/login", adminAuthHandler.Login)

	// Authenticated admin endpoints
	admin := v1.Group("")
	admin.Use(middleware.Auth(tokenSvc, tokenRepo))
	{
		admin.GET("/admin/auth/me", adminAuthHandler.Me)
		admin.POST("/admin/auth/logout", adminAuthHandler.Logout)
		admin.GET("/admin/events/:eventSlug/registrations", adminRegHandler.List)
		admin.PATCH("/admin/events/:eventSlug/registrations/:registrationId/verify", adminRegHandler.Verify)
		admin.POST("/admin/events/:eventSlug/registrations/:registrationId/resend-ticket", adminRegHandler.ResendTicket)
		admin.POST("/admin/events/:eventSlug/registrations/resend-all-tickets", adminRegHandler.ResendAllTickets)
	}

	return r
}
