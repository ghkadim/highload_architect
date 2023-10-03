BINARY_APP_NAME=bin/app
BINARY_DIALOG_NAME=bin/dialog
BINARY_RESHARD_NAME=bin/reshard
COMPOSE="docker-compose.yml"

.PHONY: build
build: test
	GOARCH=amd64 go build -o ${BINARY_APP_NAME} cmd/app/main.go
	GOARCH=amd64 go build -o ${BINARY_DIALOG_NAME} cmd/dialog/main.go
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
	docker-compose -f docker-compose.yml \
		-f docker-compose_test.yml \
		up --abort-on-container-exit --exit-code-from test

	echo Test HA setup
	docker-compose -f docker-compose_ha.yml \
		-f docker-compose_test.yml \
		up --abort-on-container-exit --exit-code-from test

#	echo Test dialogs in mysql
#	IN_MEMORY_DIALOG_ENABLED=false \
#	docker-compose -f docker-compose.yml \
#		-f docker-compose_test.yml \
#		up --abort-on-container-exit --exit-code-from test

#	echo Test dialogs sharding
#	IN_MEMORY_DIALOG_ENABLED=false \
#	docker-compose -f docker-compose.yml \
#		-f docker-compose_sharding.yml \
#		-f docker-compose_test.yml \
#		up --abort-on-container-exit --exit-code-from test

#	echo Test with cache disabled
#	CACHE_ENABLED=false \
#	docker-compose -f docker-compose.yml \
#		-f docker-compose_test.yml \
#		up --abort-on-container-exit --exit-code-from test

.PHONY: generate
generate:
	for service in app dialog ; do \
  		mkdir -p generated/$$service/go_server; \
  		cp generated/go_server-openapi-generator-ignore generated/$$service/go_server/.openapi-generator-ignore; \
		openapi-generator generate \
			-i api/$$service/openapi.json \
			-g go-server \
			-o generated/$$service/go_server || exit 1; \
		if [ "$$service" == "app" ]; then \
			PACKAGE_NAME="openapi_client"; \
		else \
			PACKAGE_NAME="openapi_client_$$service"; \
		fi; \
		openapi-generator generate \
			-i api/$$service/openapi.json \
			-g python-prior \
			-o generated/$$service/python_client \
			--package-name=$$PACKAGE_NAME || exit 1; \
		python3 generated/patch_go_server.py ./api/$$service/openapi.json \
			| gofmt | tee generated/$$service/go_server/go/authorize_routes.go; \
	done
	protoc --go_out=${GOPATH}/src --go_opt=paths=import \
    	--go-grpc_out=${GOPATH}/src --go-grpc_opt=paths=import \
    	api/dialog/dialog.proto
	go generate ./...
