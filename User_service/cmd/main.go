package main

import (
	"fmt"
	"log"

	"github.com/baigel/lms/user-service/internal/config"
	"github.com/baigel/lms/user-service/internal/handler"
	"github.com/baigel/lms/user-service/internal/middleware"
	"github.com/baigel/lms/user-service/internal/service"
	"github.com/baigel/lms/user-service/keycloak"
	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()
	if cfg.AppPort == "" {
		cfg.AppPort = "8081"
	}

	// Сборка зависимостей (dependency injection вручную)
	kc := keycloak.New(cfg.KeycloakURL, cfg.KeycloakClientID, cfg.KeycloakSecret, cfg.KeycloakRealm)
	svc := service.New(kc)
	h := handler.New(svc)
	mw := middleware.New(kc)

	router := gin.Default()
	auth := router.Group("/auth")
	{
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.Refresh)
		auth.POST("/logout", h.Logout)
		auth.POST("/register", mw.RequireAdmin(), h.Register)
	}

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("User Service is running on port %s", addr)
	log.Fatal(router.Run(addr))
}
