package usecase

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/PLoOc/case/aggregator/internal/domain"
	"github.com/PLoOc/case/aggregator/internal/infra/queue"
	"github.com/PLoOc/case/aggregator/internal/infra/repository"
)

type AggregateUseCase struct {
	consumer *queue.SQSConsumer
	repo     *repository.DynamoDBRepository
}

func NewAggregateUseCase(consumer *queue.SQSConsumer, repo *repository.DynamoDBRepository) *AggregateUseCase {
	return &AggregateUseCase{consumer: consumer, repo: repo}
}

func (uc *AggregateUseCase) Run(ctx context.Context) {
	log.Println("Aggregator consumindo mensagens...")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			msgs, err := uc.consumer.Receive(ctx)
			if err != nil {
				log.Printf("Erro ao receber mensagem: %v\n", err)
				continue
			}
			for _, msg := range msgs.Messages {
				uc.process(ctx, *msg.Body, *msg.ReceiptHandle)
			}
		}
	}
}

func (uc *AggregateUseCase) process(ctx context.Context, body, receipt string) {
	var event domain.EventRecord
	if err := json.Unmarshal([]byte(body), &event); err != nil {
		log.Printf("Erro ao fazer parse: %v\n", err)
		return
	}

	// Idempotência: ignora duplicados
	exists, err := uc.repo.EventExists(ctx, event.EventID)
	if err != nil {
		log.Printf("Erro ao verificar idempotência: %v\n", err)
		return
	}
	if exists {
		log.Printf("Evento duplicado ignorado: %s\n", event.EventID)
		uc.consumer.Delete(ctx, receipt)
		return
	}

	event.ProcessedAt = time.Now()

	if err := uc.repo.SaveEvent(ctx, &event); err != nil {
		log.Printf("Erro ao salvar evento: %v\n", err)
		return
	}

	if err := uc.repo.UpdateSummary(ctx, &event); err != nil {
		log.Printf("Erro ao atualizar summary: %v\n", err)
		return
	}

	uc.consumer.Delete(ctx, receipt)
	log.Printf("Evento agregado com sucesso: %s\n", event.EventID)
}