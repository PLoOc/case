FROM golang:1.21-alpine as builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o aggregator ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/aggregator .

CMD ["./aggregator"]