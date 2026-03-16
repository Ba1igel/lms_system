package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/baigel/lms/user-service/internal/config"
	"github.com/baigel/lms/user-service/internal/handler"
)

func main() {
	fmt.Println("User service has been started")
	cfg := config.Load()

	if cfg.AppPort == "" {
		cfg.AppPort = "8081" // default port
	}

	h := handler.New(cfg)

	http.HandleFunc("/auth/register", h.Register)
	http.HandleFunc("/auth/login", h.Login)
	http.HandleFunc("/auth/refresh", h.Refresh)
	http.HandleFunc("/auth/logout", h.Logout)

	addr := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("User Service is running on port %s", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
