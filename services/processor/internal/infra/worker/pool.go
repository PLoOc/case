package worker

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/PLoOc/case/processor/internal/domain"
	"github.com/PLoOc/case/processor/internal/infra/queue"
)

type WorkerPool struct {
	sqs    *queue.SQSAdapter
	workers int
	processorID string
}

func NewWorkerPool(sqs *queue.SQSAdapter, workers int, processorID string) *WorkerPool {
	return &WorkerPool{
		sqs:     sqs,
		workers: workers,
		processorID: processorID,
	}
}

func (wp *WorkerPool) Start(ctx context.Context) {
	var wg sync.WaitGroup
	
	for i := 0; i < wp.workers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			wp.worker(ctx, workerID)
		}(i)
	}

	<-ctx.Done()
	wg.Wait()
	log.Println("Worker pool finalizado")
}

func (wp *WorkerPool) worker(ctx context.Context, id int) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			messages, err := wp.sqs.ReceiveMessage(ctx)
			if err != nil {
				log.Printf("Worker %d: erro ao receber mensagem: %v\n", id, err)
				continue
			}

			for _, msg := range messages.Messages {
				wp.processMessage(ctx, msg)
			}
		}
	}
}

func (wp *WorkerPool) processMessage(ctx context.Context, msg types.Message) {

	var rawEvent domain.RawEvent
	
	err := json.Unmarshal([]byte(*msg.Body), &rawEvent)
	if err != nil {
		log.Printf("Erro ao fazer parse do evento: %v\n", err)
		wp.sqs.DeleteMessage(ctx, *msg.ReceiptHandle)
		return
	}

	// Validar e converter
	processed, err := rawEvent.ValidateAndConvert(wp.processorID)
	if err != nil {
		log.Printf("Evento inválido [%s]: %v\n", rawEvent.EventID, err)
		//3 tentativas de reprocessamento
		return
	}

	// Publicar na fila de eventos processados
	err = wp.sqs.SendProcessedEvent(ctx, processed)
	if err != nil {
		log.Printf("Erro ao publicar evento processado: %v\n", err)
		return
	}

	wp.sqs.DeleteMessage(ctx, *msg.ReceiptHandle)
	log.Printf("Evento processado com sucesso: %s\n", rawEvent.EventID)
}