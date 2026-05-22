package domain

import "time"

type EventRecord struct {
	EventID     string    `json:"event_id"`
	DeveloperID string    `json:"developer_id"`
	MetricType  string    `json:"metric_type"`
	Value       int       `json:"value"`
	Repository  string    `json:"repository"`
	Timestamp   time.Time `json:"timestamp"`
	ProcessedAt time.Time `json:"processed_at"`
}

type DeveloperSummary struct {
	DeveloperID        string    `json:"developer_id"`
	TotalCommits       int       `json:"total_commits"`
	TotalPullRequests  int       `json:"total_pull_requests"`
	AvgReviewTime      float64   `json:"avg_review_time_minutes"`
	EventsProcessed    int       `json:"events_processed"`
	LastActivity       time.Time `json:"last_activity"`
}