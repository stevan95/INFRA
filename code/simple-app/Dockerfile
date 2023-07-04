##
FROM golang:1.16-alpine AS builder

WORKDIR app
COPY app/ .

RUN go mod download
RUN go mod verify
RUN go build cmd/web/*

FROM alpine:3.17

WORKDIR app

RUN mkdir templates
COPY app/templates/ ./templates/
COPY --from=builder /go/app/main .

EXPOSE 8080

USER 1000:1000

ENTRYPOINT ["./main"]
