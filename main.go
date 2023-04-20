package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/ghkadim/highload_architect/internal/cache"
	"github.com/ghkadim/highload_architect/internal/controller"
	"github.com/ghkadim/highload_architect/internal/logger"
	"github.com/ghkadim/highload_architect/internal/mysql"
	"github.com/ghkadim/highload_architect/internal/session"

	"github.com/ghkadim/highload_architect/internal/service"
)

func getenv[T any](variable string, defaultValue T) T {
	valueStr := os.Getenv(variable)
	if valueStr == "" {
		return defaultValue
	}

	var value T
	_, err := fmt.Sscan(valueStr, &value)
	if err != nil {
		logger.Info("Failed to parse env variable %s value %s to type %T, return defaultValue %v",
			variable, valueStr, value, defaultValue)
		return defaultValue
	}
	return value
}

func main() {
	l := logger.Init(getenv("DEBUG", false))
	defer l.Sync()

	logger.Info("Server starting")

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
		if address == "" {
			continue
		}
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

	var cache_ service.Cache
	if getenv("CACHE_ENABLED", true) {
		cache_ = cache.NewCache(
			getenv("CACHE_FEED_LIMIT", 1000),
			cache.NewLoadWithRetry(master))
		logger.Info("Feed cache enabled")
	} else {
		cache_ = cache.NewDisabledCache()
		logger.Info("Feed cache disabled")
	}

	apiService := service.NewApiService(
		master,
		slaves,
		cache_,
		session.NewSession(
			getenv("SESSION_KEY", "secret"),
		),
	)
	apiController := controller.NewApiController(
		apiService,
	)

	log.Fatal(http.ListenAndServe(":8080", apiController))
}
