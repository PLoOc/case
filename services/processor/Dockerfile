FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go mod tidy
RUN go build -o processor ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/processor .
CMD ["./processor"]