package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/baigel/lms/user-service/internal/config"
	"github.com/baigel/lms/user-service/internal/handler"
	"github.com/baigel/lms/user-service/internal/middleware"
	"github.com/baigel/lms/user-service/internal/service"
	"github.com/baigel/lms/user-service/keycloak"
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

	http.HandleFunc("/auth/login", h.Login)

	http.HandleFunc("/auth/refresh", h.Refresh)

	http.HandleFunc("/auth/logout", h.Logout)

	http.Handle("/auth/register", mw.RequireAdmin(http.HandlerFunc(h.Register)))

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("User Service is running on port %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
