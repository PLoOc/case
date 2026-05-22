package repository

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/PLoOc/case/aggregator/internal/domain"
	"log"
)

type DynamoDBRepository struct {
	client *dynamodb.Client
}

func NewDynamoDBRepository(ctx context.Context, endpoint, region string) (*DynamoDBRepository, error) {
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	return &DynamoDBRepository{client: client}, nil
}

func (r *DynamoDBRepository) SaveEvent(ctx context.Context, event *domain.EventRecord) error {
	item, err := attributevalue.MarshalMap(event)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("events"),
		Item:      item,
	})
	return err
}

func (r *DynamoDBRepository) UpdateSummary(ctx context.Context, summary *domain.DeveloperSummary) error {
	item, err := attributevalue.MarshalMap(summary)
	if err != nil {
		return err
	}

	_, err = r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String("developer_summary"),
		Item:      item,
	})
	return err
}

func (r *DynamoDBRepository) GetSummary(ctx context.Context, developerID string) (*domain.DeveloperSummary, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String("developer_summary"),
		Key: map[string]interface{}{
			"developer_id": &dynamodb.AttributeValueMemberS{Value: developerID},
		},
	})
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return &domain.DeveloperSummary{
			DeveloperID: developerID,
		}, nil
	}

	var summary domain.DeveloperSummary
	err = attributevalue.UnmarshalMap(result.Item, &summary)
	return &summary, err
}