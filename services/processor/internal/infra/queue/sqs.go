package queue

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"log"
)

type SQSAdapter struct {
	client             *sqs.Client
	rawEventsQueueURL  string
	processedEventsQueueURL string
}

func NewSQSAdapter(ctx context.Context, endpoint, region, rawQueue, processedQueue string) (*SQSAdapter, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}

	client := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	// Obter URL da fila raw-events
	rawResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(rawQueue),
	})
	if err != nil {
		return nil, err
	}

	// Obter URL da fila processed-events
	procResp, err := client.GetQueueUrl(ctx, &sqs.GetQueueUrlInput{
		QueueName: aws.String(processedQueue),
	})
	if err != nil {
		return nil, err
	}

	return &SQSAdapter{
		client:              client,
		rawEventsQueueURL:   *rawResp.QueueUrl,
		processedEventsQueueURL: *procResp.QueueUrl,
	}, nil
}

func (s *SQSAdapter) ReceiveMessage(ctx context.Context) (*sqs.ReceiveMessageOutput, error) {
	return s.client.ReceiveMessage(ctx, &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(s.rawEventsQueueURL),
		MaxNumberOfMessages: types.Int32(10),
		WaitTimeSeconds:     10,
	})
}

func (s *SQSAdapter) DeleteMessage(ctx context.Context, receiptHandle string) error {
	_, err := s.client.DeleteMessage(ctx, &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(s.rawEventsQueueURL),
		ReceiptHandle: aws.String(receiptHandle),
	})
	return err
}

func (s *SQSAdapter) SendProcessedEvent(ctx context.Context, event interface{}) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	_, err = s.client.SendMessage(ctx, &sqs.SendMessageInput{
		QueueUrl:    aws.String(s.processedEventsQueueURL),
		MessageBody: aws.String(string(body)),
	})
	return err
}

func (s *SQSAdapter) Close() error {
	return nil // SQS não precisa de close explícito
}