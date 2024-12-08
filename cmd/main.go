package main

import (
	"astral/internal/cache"
	handler "astral/internal/contorller"
	"astral/internal/repository"
	"astral/internal/usecase"
	"astral/pkg/database"
	"astral/pkg/logger"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("HTTP_PORT")
	db := database.InitializeDBPostgres(3, 10)
	logger.InitLogger()

	repository := repository.NewRepository(db.GetDB())
	usecase := usecase.NewUseCase(repository)
	handlers := handler.NewHandler(usecase)
	cache.InitializeCache()
	router := handlers.Handle()
	err := http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatalf("Connection failed: %s\n", err.Error())
	}

	log.Infof("Server is running on port %s\n", port)
}
