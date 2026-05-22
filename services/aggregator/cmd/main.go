package main

import (
	"context"
	"log"
	"net/http"

	"github.com/PLoOc/case/aggregator/internal/infra/api"
	"github.com/PLoOc/case/aggregator/internal/infra/config"
	"github.com/PLoOc/case/aggregator/internal/infra/queue"
	"github.com/PLoOc/case/aggregator/internal/infra/repository"
	"github.com/PLoOc/case/aggregator/internal/usecase"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	repo, err := repository.NewDynamoDBRepository(ctx, cfg.AWSEndpoint, cfg.AWSRegion)
	if err != nil {
		log.Fatalf("Erro ao conectar DynamoDB: %v", err)
	}

	consumer, err := queue.NewSQSConsumer(ctx, cfg.AWSEndpoint, cfg.AWSRegion, cfg.ProcessedQueue)
	if err != nil {
		log.Fatalf("Erro ao conectar SQS: %v", err)
	}

	go usecase.NewAggregateUseCase(consumer, repo).Run(ctx)

	handler := api.NewHandler(repo)
	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/metrics/{developer_id}", handler.GetEvents)
	mux.HandleFunc("/metrics/{developer_id}/summary", handler.GetSummary)

	log.Printf("API rodando na porta %s\n", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, mux))
}