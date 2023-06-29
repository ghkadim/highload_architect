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
	"github.com/ghkadim/highload_architect/internal/models"
	"github.com/ghkadim/highload_architect/internal/mysql"
	"github.com/ghkadim/highload_architect/internal/session"

	"github.com/ghkadim/highload_architect/internal/service"
)

func main() {
	l := logger.Init(getenv("DEBUG", false))
	defer func() { _ = l.Sync() }()

	logger.Info("Server starting")

	master, err := mysql.NewStorage(
		getenv("DB_USER", "user"),
		getenv("DB_PASSWORD", "password"),
		getenv("DB_ADDRESS", "127.0.0.1:3306"),
		getenv("DB_DATABASE", "db"),
		getenv("DB_DEDICATED_SHARDS", make(mysql.DedicatedShardID)),
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
			getenv("DB_DEDICATED_SHARDS", make(mysql.DedicatedShardID)),
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

func getenv[T any](variable string, defaultValue T) T {
	valueStr := os.Getenv(variable)
	if valueStr == "" {
		return defaultValue
	}

	var value T
	err := parseValue(valueStr, &value)
	if err != nil {
		logger.Info("Failed to parse env variable %s, return defaultValue %v: %v",
			variable, defaultValue, err)
		return defaultValue
	}
	return value
}

func parseValue(valueStr string, value any) error {
	switch val := value.(type) {
	case *mysql.DedicatedShardID:
		for _, kv := range strings.Split(valueStr, ",") {
			kvArr := strings.Split(kv, ":")
			if len(kvArr) != 2 {
				return fmt.Errorf("failed to parse key value %s", kv)
			}
			var userID models.UserID
			err := parseValue(kvArr[0], &userID)
			if err != nil {
				return err
			}
			(*val)[userID] = kvArr[1]
		}
	default:
		_, err := fmt.Sscan(valueStr, val)
		if err != nil {
			return fmt.Errorf("failed to parse value '%s' to type %T: %w", valueStr, value, err)
		}
	}
	return nil
}
