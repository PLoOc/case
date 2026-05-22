package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/PLoOc/case/processor/internal/infra/queue"
	"github.com/PLoOc/case/processor/internal/infra/worker"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Carregar configurações
	endpoint := os.Getenv("AWS_ENDPOINT_URL")
	region := os.Getenv("AWS_REGION")
	rawQueue := os.Getenv("RAW_EVENTS_QUEUE")
	processedQueue := os.Getenv("PROCESSED_EVENTS_QUEUE")
	workerCount, _ := strconv.Atoi(os.Getenv("WORKER_COUNT"))

	if endpoint == "" {
		endpoint = "http://localhost:4566"
	}
	if region == "" {
		region = "us-east-1"
	}
	if rawQueue == "" {
		rawQueue = "raw-events"
	}
	if processedQueue == "" {
		processedQueue = "processed-events"
	}
	if workerCount == 0 {
		workerCount = 3
	}

	log.Println("🚀 Processor iniciando...")
	log.Printf("  Endpoint: %s\n", endpoint)
	log.Printf("  Workers: %d\n", workerCount)

	// Conectar ao SQS
	sqs, err := queue.NewSQSAdapter(ctx, endpoint, region, rawQueue, processedQueue)
	if err != nil {
		log.Fatalf("Erro ao conectar ao SQS: %v", err)
	}
	defer sqs.Close()

	// Iniciar worker pool
	pool := worker.NewWorkerPool(sqs, workerCount, "processor-1")
	pool.Start(ctx)
}