package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	domain "github.com/PLoOc/case/aggregator/internal/domain"
)

type DynamoDBRepository struct {
	client *dynamodb.Client
}

func NewDynamoDBRepository(ctx context.Context, endpoint, region string) (*DynamoDBRepository, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(aws.CredentialsProviderFunc(func(ctx context.Context) (aws.Credentials, error) {
			return aws.Credentials{AccessKeyID: "test", SecretAccessKey: "test"}, nil
		})),
	)
	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	return &DynamoDBRepository{client: client}, nil
}

func (r *DynamoDBRepository) EventExists(ctx context.Context, eventID string) (bool, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String("events"),
		Key: map[string]types.AttributeValue{
			"event_id": &types.AttributeValueMemberS{Value: eventID},
		},
	})
	if err != nil {
		return false, err
	}
	return result.Item != nil, nil
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

func (r *DynamoDBRepository) GetSummary(ctx context.Context, developerID string) (*domain.DeveloperSummary, error) {
	result, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String("developer_summary"),
		Key: map[string]types.AttributeValue{
			"developer_id": &types.AttributeValueMemberS{Value: developerID},
		},
	})
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return &domain.DeveloperSummary{DeveloperID: developerID}, nil
	}
	var summary domain.DeveloperSummary
	err = attributevalue.UnmarshalMap(result.Item, &summary)
	return &summary, err
}

func (r *DynamoDBRepository) UpdateSummary(ctx context.Context, event *domain.EventRecord) error {
	summary, err := r.GetSummary(ctx, event.DeveloperID)
	if err != nil {
		return err
	}

	summary.EventsProcessed++
	summary.LastActivity = event.Timestamp

	switch event.MetricType {
	case "commits":
		summary.TotalCommits += event.Value
	case "pull_requests":
		summary.TotalPullRequests += event.Value
	case "review_time_minutes":
		summary.TotalReviewTime += event.Value
		if summary.EventsProcessed > 0 {
			summary.AvgReviewTime = float64(summary.TotalReviewTime) / float64(summary.EventsProcessed)
		}
	default:
		return fmt.Errorf("metric_type desconhecido: %s", event.MetricType)
	}

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

func (r *DynamoDBRepository) GetEvents(ctx context.Context, developerID string) ([]domain.EventRecord, error) {
	result, err := r.client.Scan(ctx, &dynamodb.ScanInput{
		TableName:        aws.String("events"),
		FilterExpression: aws.String("developer_id = :dev"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":dev": &types.AttributeValueMemberS{Value: developerID},
		},
	})
	if err != nil {
		return nil, err
	}

	var events []domain.EventRecord
	err = attributevalue.UnmarshalListOfMaps(result.Items, &events)
	return events, err
}