FROM golang:1.20 AS build
WORKDIR /src
COPY . /src

ENV CGO_ENABLED=0
RUN go get -d -v ./...

RUN make build

FROM scratch AS runtime
ARG APPLICATION=app
COPY --from=build /src/bin/$APPLICATION ./app
EXPOSE 8080/tcp
ENTRYPOINT ["./app"]
