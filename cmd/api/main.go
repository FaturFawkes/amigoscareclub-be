package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"myapp/application/usecase"
	httpDelivery "myapp/delivery/http"
	"myapp/infrastructure/auth"
	"myapp/infrastructure/config"
	infraDB "myapp/infrastructure/db"
	"myapp/infrastructure/idgen"
	"myapp/infrastructure/repository/postgres"
	infraS3 "myapp/infrastructure/s3"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	ctx := context.Background()
	pool, err := infraDB.NewPool(ctx, cfg)
	if err != nil {
		log.Fatalf("database: %v", err)
	}
	defer pool.Close()

	// Repositories
	regRepo := postgres.NewRegistrationRepo(pool)
	eventRepo := postgres.NewEventRepo(pool)
	adminRepo := postgres.NewAdminRepo(pool)
	tokenRepo := postgres.NewRevokedTokenRepo(pool)

	// Services
	jwtSvc := auth.NewJWTService(cfg)
	bcryptHasher := auth.NewBcryptHasher()
	idGen := idgen.New(pool, cfg)
	storage, err := infraS3.New(cfg)
	if err != nil {
		log.Fatalf("s3: %v", err)
	}

	realClock := &realClock{}

	// Use cases — public
	createRegUC := usecase.NewCreateRegistrationUseCase(regRepo, eventRepo, storage, nil, idGen, realClock)
	getRegUC := usecase.NewGetRegistrationUseCase(regRepo)
	getEventUC := usecase.NewGetEventUseCase(eventRepo)

	// Use cases — admin
	loginUC := usecase.NewAdminLoginUseCase(adminRepo, bcryptHasher, jwtSvc)
	meUC := usecase.NewAdminMeUseCase(adminRepo)
	logoutUC := usecase.NewAdminLogoutUseCase(tokenRepo, jwtSvc)
	listRegsUC := usecase.NewAdminListRegistrationsUseCase(regRepo)
	verifyRegUC := usecase.NewAdminVerifyRegistrationUseCase(regRepo, realClock)

	// Handlers
	regHandler := httpDelivery.NewRegistrationHandler(createRegUC, getRegUC)
	eventHandler := httpDelivery.NewEventHandler(getEventUC)
	adminAuthHandler := httpDelivery.NewAdminAuthHandler(loginUC, meUC, logoutUC)
	adminRegHandler := httpDelivery.NewAdminRegistrationHandler(listRegsUC, verifyRegUC)

	router := httpDelivery.NewRouter(cfg.APIBasePath, jwtSvc, tokenRepo, regHandler, eventHandler, adminAuthHandler, adminRegHandler)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		log.Printf("server listening on :%s (basePath=%s)", cfg.Port, cfg.APIBasePath)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")
	shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutCtx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}
	log.Println("server exited")
}

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }
