#
FROM golang:1.16-alpine AS builder

WORKDIR app
COPY app/ .

RUN go mod download
RUN go mod verify
RUN go build cmd/web/*

EXPOSE 8080

ENTRYPOINT ["./main"]
