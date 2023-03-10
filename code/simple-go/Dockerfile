FROM golang:1.16-alpine AS build

WORKDIR app

COPY app/ .

RUN go build cmd/web/*

EXPOSE 8080

USER 1000:1000

ENTRYPOINT ["./main"]
