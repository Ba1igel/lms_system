package app

import (
	"URL_SHORTENER/internal/handler"
	"URL_SHORTENER/internal/repo"
	"URL_SHORTENER/internal/service"
	"log"
	"net/http"
)

func main() {
	cfg := config.InitConfig()
	r := repo.NewUrlRepo(db)
	s := service.NewUrlService(r)
	urlHandler := handler.NewUrlHandler(s)

	http.HandleFunc("/shorten", urlHandler.code)
	log.Println("Starting server on localhost:8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("error starting server: %v", err)
	}
}
