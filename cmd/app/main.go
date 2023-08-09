package main

import (
	"net/http"
	"strings"

	"github.com/ghkadim/highload_architect/internal/app/cache"
	"github.com/ghkadim/highload_architect/internal/app/controller/openapi"
	"github.com/ghkadim/highload_architect/internal/app/controller/websocket"
	dialogSvc "github.com/ghkadim/highload_architect/internal/app/dialog"
	"github.com/ghkadim/highload_architect/internal/app/mysql"
	"github.com/ghkadim/highload_architect/internal/app/rabbitmq"
	"github.com/ghkadim/highload_architect/internal/app/service/dialog"
	"github.com/ghkadim/highload_architect/internal/app/service/friend"
	"github.com/ghkadim/highload_architect/internal/app/service/post"
	"github.com/ghkadim/highload_architect/internal/app/service/user"
	"github.com/ghkadim/highload_architect/internal/config"
	"github.com/ghkadim/highload_architect/internal/dialog/tarantool"
	"github.com/ghkadim/highload_architect/internal/logger"
	"github.com/ghkadim/highload_architect/internal/server"
	"github.com/ghkadim/highload_architect/internal/session"
)

func main() {
	l := logger.Init(config.Get("DEBUG", false))
	defer func() { _ = l.Sync() }()

	logger.Info("Server starting")

	storage, err := mysql.NewStorage(
		config.Get("DB_USER", "user"),
		config.Get("DB_PASSWORD", "password"),
		config.Get("DB_ADDRESS", "127.0.0.1:3306"),
		config.Get("DB_DATABASE", "db"),
		config.Get("DB_DEDICATED_SHARDS", make(mysql.DedicatedShardID)),
	)
	if err != nil {
		logger.Fatal("failed to init db")
	}

	replicas := make([]mysql.Storage, 0)
	for _, address := range strings.Split(config.Get("DB_REPLICA_ADDRESSES", ""), ",") {
		if address == "" {
			continue
		}
		replica, err := mysql.NewStorage(
			config.Get("DB_USER", "user"),
			config.Get("DB_PASSWORD", "password"),
			address,
			config.Get("DB_DATABASE", "db"),
			config.Get("DB_DEDICATED_SHARDS", make(mysql.DedicatedShardID)),
		)
		if err != nil {
			logger.Fatal("failed to init slave db")
		}

		replicas = append(replicas, replica)
	}

	if len(replicas) > 0 {
		storage = mysql.NewStorageWithReplicas(storage, replicas)
	}

	var cache_ cache.Cache
	if config.Get("CACHE_ENABLED", true) {
		cache_ = cache.NewCache(
			config.Get("CACHE_FEED_LIMIT", 1000),
			cache.NewLoadWithRetry(storage))
		logger.Info("Feed cache enabled")
	} else {
		cache_ = cache.NewDisabledCache()
		logger.Info("Feed cache disabled")
	}

	session_ := session.NewSession(
		config.Get("SESSION_KEY", "secret"),
	)

	var eventConsumer rabbitmq.Consumer
	var eventPublisher rabbitmq.Publisher

	if config.Get("ASYNCAPI_ENABLED", true) {
		eventConsumer, err = rabbitmq.NewConsumer(
			config.Get("RMQ_USER", "guest"),
			config.Get("RMQ_PASSWORD", "guest"),
			config.Get("RMQ_ADDRESS", "localhost:5672"),
			config.Get("EVENT_CONSUMER_QUEUE_LEN", 1000),
		)
		if err != nil {
			logger.Fatal("Failed to create event consumer: %v", err)
		}

		eventPublisher, err = rabbitmq.NewPublisher(
			config.Get("RMQ_USER", "guest"),
			config.Get("RMQ_PASSWORD", "guest"),
			config.Get("RMQ_ADDRESS", "localhost:5672"),
		)
		if err != nil {
			logger.Fatal("Failed to create event publisher: %v", err)
		}
	} else {
		eventConsumer = rabbitmq.NewNopConsumer()
		eventPublisher = rabbitmq.NewNopPublisher()
	}

	var dialogService dialog.Service
	if config.Get("DIALOG_MICROSERVICE_ENABLED", false) {
		dialogService, err = dialogSvc.NewClient(config.Get("DIALOG_ADDRESS", "localhost:8081"))
		if err != nil {
			logger.Fatal("Failed to create client for dialog service: %v", err)
		}
	} else {
		if config.Get("IN_MEMORY_DIALOG_ENABLED", true) {
			logger.Info("In memory dialogs enabled")
			dialogService = dialog.NewService(
				tarantool.NewStorage(
					config.Get("TARANTOOL_ADDRESS", ""),
					http.Client{},
				),
			)
		} else {
			logger.Info("In memory dialogs disabled")
			dialogService = dialog.NewService(storage)
		}
	}

	apiService := &svc{
		userService:   user.NewService(storage, session_),
		friendService: friend.NewService(storage, cache_, eventPublisher),
		postService:   post.NewService(storage, cache_, eventPublisher, eventConsumer),
		dialogService: dialogService,
	}

	routers := []server.Router{
		openapi.NewRouter(
			openapi.NewController(apiService, session_)),
	}

	if config.Get("ASYNCAPI_ENABLED", true) {
		logger.Info("Async API enabled")
		routers = append(routers, websocket.NewRouter(
			websocket.NewController(apiService, session_)))
	}

	srv := server.NewServer(routers...)
	srv.ListenAndServe(":8080")
}

type userService = *user.Service
type friendService = *friend.Service
type postService = *post.Service
type dialogService = dialog.Service

type svc struct {
	userService
	friendService
	postService
	dialogService
}
