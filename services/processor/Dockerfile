FROM golang:1.21-alpine as builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o processor ./cmd/main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/processor .

CMD ["./processor"]