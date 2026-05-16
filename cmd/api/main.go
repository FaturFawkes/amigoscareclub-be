package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"

	"myapp/application/usecase"
	httpDelivery "myapp/delivery/http"
	"myapp/infrastructure/repository/mysql"
)

func main() {
	repo := mysql.NewRegistrationRepository(nil)
	placeUseCase := usecase.NewPlaceRegistrationUseCase(repo, nil, nil)
	getUseCase := usecase.NewGetRegistrationUseCase(repo)
	
	handler := httpDelivery.NewRegistrationHandler(placeUseCase, getUseCase)

	router := gin.Default()
	router.POST("/registrations", func(c *gin.Context) {
		handler.Place(c.Writer, c.Request)
	})
	router.GET("/registrations", func(c *gin.Context) {
		handler.Get(c.Writer, c.Request)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
