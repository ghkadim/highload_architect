BINARY_NAME=bin/app

build:
	GOARCH=amd64 go build -o ${BINARY_NAME} main.go

run:
	./${BINARY_NAME}

build_and_run: build run

clean:
	go clean
	rm ${BINARY_NAME}

compose_build:
	docker-compose build

compose_run:
	docker-compose up

compose_build_and_run: compose_build compose_run

compose_clean:
	docker-compose down
	docker volume rm highload_architect_my-db

generate:
	openapi-generator generate \
		-i ./api/openapi.json \
		-g go-server \
		-o ./generated
