FROM golang:1.16-alpine AS builder

WORKDIR app
COPY app/ .

RUN go mod download
RUN go mod verify
RUN go build 

FROM alpine:3.17

WORKDIR app

RUN mkdir templates && mkdir cert
COPY app/templates/ ./templates/
COPY app/cert ./cert/
COPY --from=builder /go/app/imageTagAdmissionController .

EXPOSE 8443

USER 1000:1000

ENTRYPOINT ["./imageTagAdmissionController.exe"]