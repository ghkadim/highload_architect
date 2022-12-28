BINARY_NAME=bin/app

build:
	GOARCH=amd64 go build -o ${BINARY_NAME} main.go

run:
	./${BINARY_NAME}

build_and_run: build run

clean:
	go clean
	rm ${BINARY_NAME}

generate:
	openapi-generator generate \
		-i ./api/openapi.json \
		-g go-server \
		-o ./generated
