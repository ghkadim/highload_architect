BINARY_APP_NAME=bin/app
BINARY_RESHARD_NAME=bin/reshard
COMPOSE="docker-compose.yml"

.PHONY: build
build: test
	GOARCH=amd64 go build -o ${BINARY_APP_NAME} cmd/app/main.go
	GOARCH=amd64 go build -o ${BINARY_RESHARD_NAME} cmd/reshard/main.go

.PHONY: run
run:
	./${BINARY_APP_NAME}

.PHONY: build_and_run
build_and_run: build run

.PHONY: test
test:
	go test ./...

.PHONY: linter-install
linter-install:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.53.3

.PHONY: lint
lint:
	golangci-lint run --timeout 5m --verbose

.PHONY: test_api
test_api:
	pytest -x -v test/e2e

.PHONY: clean
clean:
	go clean
	rm ${BINARY_APP_NAME}

.PHONY: compose_build
compose_build:
	docker-compose -f ${COMPOSE} build

.PHONY: compose_run
compose_run:
	docker-compose -f ${COMPOSE} up

.PHONY: compose_build_and_run
compose_build_and_run: compose_build compose_run

.PHONY: compose_clean
compose_clean:
	docker-compose -f ${COMPOSE} down
	docker volume ls -q | grep highload_architect | xargs docker volume rm

.PHONY: compose_test
compose_test: compose_build
	echo Test dialogs in tarantool
	docker-compose -f docker-compose.yml \
		-f docker-compose_test.yml \
		up --abort-on-container-exit --exit-code-from test

	echo Test dialogs in mysql
	IN_MEMORY_DIALOG_ENABLED=false \
	docker-compose -f docker-compose.yml \
		-f docker-compose_test.yml \
		up --abort-on-container-exit --exit-code-from test

	echo Test dialogs sharding
	IN_MEMORY_DIALOG_ENABLED=false \
	docker-compose -f docker-compose.yml \
		-f docker-compose_sharding.yml \
		-f docker-compose_test.yml \
		up --abort-on-container-exit --exit-code-from test

	echo Test with cache disabled
	CACHE_ENABLED=false \
	docker-compose -f docker-compose.yml \
		-f docker-compose_test.yml \
		up --abort-on-container-exit --exit-code-from test

.PHONY: generate
generate:
	openapi-generator generate \
		-i ./api/openapi.json \
		-g go-server \
		-o ./generated/go_server
	openapi-generator generate \
		-i ./api/openapi.json \
		-g python-prior \
		-o ./generated/python_client
	python3 generated/patch_go_server.py | gofmt | tee "generated/go_server/go/authorize_routes.go"
	go generate ./...
