FROM golang:1.16.5 AS development
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go install github.com/cespare/reflex@latest
EXPOSE 4000
CMD reflex -g '*.go' go run main.go --start-service

FROM golang:1.16.5 AS builder
ENV GOOS linux
ENV CGO_ENABLED 0
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o app

FROM alpine:latest AS production
RUN apk add --no-cache ca-certificates
COPY --from=builder app .
EXPOSE 4000
CMD ./app
