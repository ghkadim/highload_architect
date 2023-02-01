FROM golang:1.19 AS build
WORKDIR /src
COPY . /src

ENV CGO_ENABLED=0
RUN go get -d -v ./...

RUN make build

FROM scratch AS runtime
COPY --from=build /src/bin/app ./
EXPOSE 8080/tcp
ENTRYPOINT ["./app"]
