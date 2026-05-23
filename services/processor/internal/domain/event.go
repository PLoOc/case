package domain

import (
	"errors"
	"fmt"
	"regexp"
	"time"
)

type RawEvent struct {
	EventID    string    `json:"event_id"`
	DeveloperID string   `json:"developer_id"`
	MetricType string    `json:"metric_type"`
	Value      int       `json:"value"`
	Repository string    `json:"repository"`
	Timestamp  time.Time `json:"timestamp"`
}

type ProcessedEvent struct {
	EventID     string    `json:"event_id"`
	DeveloperID string    `json:"developer_id"`
	MetricType  string    `json:"metric_type"`
	Value       int       `json:"value"`
	Repository  string    `json:"repository"`
	Timestamp   time.Time `json:"timestamp"`
	ProcessedAt time.Time `json:"processed_at"`
	ProcessorID string    `json:"processor_id"`
}

// ValidateAndConvert valida o evento e o converte para ProcessedEvent
func (e *RawEvent) ValidateAndConvert(processorID string) (*ProcessedEvent, error) {
	// Validar event_id (UUID)
	if e.EventID == "" {
		return nil, errors.New("event_id é obrigatório")
	}

	var uuidRegex = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$`)

	if !uuidRegex.MatchString(e.EventID) {
    	return nil, errors.New("event_id deve ser um UUID válido")
}

	// Validar developer_id
	if e.DeveloperID == "" {
		return nil, errors.New("developer_id é obrigatório")
	}

	// Validar metric_type
	validTypes := map[string]bool{
		"commits":            true,
		"pull_requests":      true,
		"review_time_minutes": true,
	}
	if !validTypes[e.MetricType] {
		return nil, fmt.Errorf("metric_type inválido: %s", e.MetricType)
	}

	// Validar value
	if e.Value < 0 {
		return nil, errors.New("value não pode ser negativo")
	}
	if e.MetricType == "review_time_minutes" && e.Value > 1440 {
		return nil, errors.New("review_time_minutes não pode ser maior que 1440")
	}

	// Validar timestamp (não pode ser no futuro)
	if e.Timestamp.After(time.Now()) {
		return nil, errors.New("timestamp não pode ser no futuro")
	}

	// Convertendo para ProcessedEvent
	return &ProcessedEvent{
		EventID:     e.EventID,
		DeveloperID: e.DeveloperID,
		MetricType:  e.MetricType,
		Value:       e.Value,
		Repository:  e.Repository,
		Timestamp:   e.Timestamp,
		ProcessedAt: time.Now(),
		ProcessorID: processorID,
	}, nil
}