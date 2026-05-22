FROM golang:1.21-alpine

WORKDIR /app

COPY ./services/processor/go.mod ./services/processor/
COPY ./services/aggregator/go.mod ./services/aggregator/

RUN mkdir -p /app/services/processor/cmd /app/services/processor/internal
RUN mkdir -p /app/services/aggregator/cmd /app/services/aggregator/internal

WORKDIR /app/services/processor
RUN go mod tidy

WORKDIR /app/services/aggregator
RUN go mod tidy

WORKDIR /app