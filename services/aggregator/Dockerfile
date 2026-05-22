FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go mod tidy
RUN go build -o aggregator ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/aggregator .
CMD ["./aggregator"]