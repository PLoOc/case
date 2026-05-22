package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/PLoOc/case/aggregator/internal/infra/api"
	"github.com/PLoOc/case/aggregator/internal/infra/repository"
)

func main() {
	ctx := context.Background()

	endpoint := os.Getenv("AWS_ENDPOINT_URL")
	if endpoint == "" {
		endpoint = "http://localhost:4566"
	}

	repo, err := repository.NewDynamoDBRepository(ctx, endpoint, "us-east-1")
	if err != nil {
		log.Fatalf("Erro ao conectar ao DynamoDB: %v", err)
	}

	handler := api.NewAPIHandler(repo)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", handler.Health)
	mux.HandleFunc("/metrics/{developer_id}/summary", handler.GetSummary)

	log.Println("🚀 Aggregator rodando em :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}