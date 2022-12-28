package main

import (
	"github.com/ghkadim/highload_architect/internal/mysql"
	"github.com/ghkadim/highload_architect/internal/session"
	"log"
	"net/http"
	"os"

	openapi "github.com/ghkadim/highload_architect/generated/go"
	"github.com/ghkadim/highload_architect/internal/service"
)

func getenv(variable, defaultValue string) string {
	value := os.Getenv(variable)
	if value == "" {
		return defaultValue
	}
	return value
}

func main() {
	log.Print("Server started")

	storage, err := mysql.NewStorage(
		getenv("DB_USER", "user"),
		getenv("DB_PASSWORD", "password"),
		getenv("DB_ADDRESS", "127.0.0.1:3306"),
		getenv("DB_DATABASE", "db"),
	)
	if err != nil {
		log.Fatal("failed to init db")
	}

	DefaultApiService := service.NewApiService(
		storage,
		session.NewSession(
			getenv("SESSION_KEY", "secret"),
		),
	)
	DefaultApiController := openapi.NewDefaultApiController(
		DefaultApiService,
	)

	router := openapi.NewRouter(DefaultApiController)

	log.Fatal(http.ListenAndServe(":8080", router))
}
