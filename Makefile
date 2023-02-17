BINARY_NAME=bin/app
COMPOSE="docker-compose.yaml"

build:
	GOARCH=amd64 go build -o ${BINARY_NAME} main.go

run:
	./${BINARY_NAME}

build_and_run: build run

clean:
	go clean
	rm ${BINARY_NAME}

compose_build:
	docker-compose -f ${COMPOSE} build

compose_run:
	docker-compose -f ${COMPOSE} up

compose_build_and_run: compose_build compose_run

compose_clean:
	docker-compose -f ${COMPOSE} down
	docker volume ls -q | grep highload_architect | xargs docker volume rm

generate:
	openapi-generator generate \
		-i ./api/openapi.json \
		-g go-server \
		-o ./generated
