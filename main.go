package main

import (
	"github.com/ghkadim/highload_architect/internal/controller"
	"github.com/ghkadim/highload_architect/internal/mysql"
	"github.com/ghkadim/highload_architect/internal/session"
	"log"
	"net/http"
	"os"
	"strings"

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

	master, err := mysql.NewStorage(
		getenv("DB_USER", "user"),
		getenv("DB_PASSWORD", "password"),
		getenv("DB_ADDRESS", "127.0.0.1:3306"),
		getenv("DB_DATABASE", "db"),
	)
	if err != nil {
		log.Fatal("failed to init db")
	}

	slaves := make([]service.Storage, 0)
	for _, address := range strings.Split(getenv("DB_REPLICA_ADDRESSES", ""), ",") {
		slave, err := mysql.NewStorage(
			getenv("DB_USER", "user"),
			getenv("DB_PASSWORD", "password"),
			address,
			getenv("DB_DATABASE", "db"),
		)
		if err != nil {
			log.Fatal("failed to init slave db")
		}

		slaves = append(slaves, slave)
	}

	apiService := service.NewApiService(
		master,
		slaves,
		session.NewSession(
			getenv("SESSION_KEY", "secret"),
		),
	)
	apiController := controller.NewApiController(
		apiService,
	)

	log.Fatal(http.ListenAndServe(":8080", apiController))
}
