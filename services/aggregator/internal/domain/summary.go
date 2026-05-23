package domain

import "time"

type EventRecord struct {
	EventID     string    `json:"event_id" dynamodbav:"event_id"`
	DeveloperID string    `json:"developer_id" dynamodbav:"developer_id"`
	MetricType  string    `json:"metric_type" dynamodbav:"metric_type"`
	Value       int       `json:"value" dynamodbav:"value"`
	Repository  string    `json:"repository" dynamodbav:"repository"`
	Timestamp   time.Time `json:"timestamp" dynamodbav:"timestamp"`
	ProcessedAt time.Time `json:"processed_at" dynamodbav:"processed_at"`
}

type DeveloperSummary struct {
	DeveloperID      string    `json:"developer_id" dynamodbav:"developer_id"`
	TotalCommits      int       `json:"total_commits" dynamodbav:"total_commits"`
	TotalPullRequests int       `json:"total_pull_requests" dynamodbav:"total_pull_requests"`
	TotalReviewTime   int       `json:"total_review_time" dynamodbav:"total_review_time"`
	AvgReviewTime     float64   `json:"avg_review_time_minutes" dynamodbav:"avg_review_time_minutes"`
	EventsProcessed   int       `json:"events_processed" dynamodbav:"events_processed"`
	LastActivity      time.Time `json:"last_activity" dynamodbav:"last_activity"`
}