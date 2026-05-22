#!/bin/bash

echo "Enviando eventos de teste..."

# Evento válido 1
aws --endpoint-url=http://localhost:4566 sqs send-message \
  --queue-url http://localhost:4566/000000000000/raw-events \
  --message-body '{
    "event_id": "550e8400-e29b-41d4-a716-446655440001",
    "developer_id": "dev-123",
    "metric_type": "commits",
    "value": 5,
    "repository": "org/repo-1",
    "timestamp": "2026-05-21T10:00:00Z"
  }'

# Mais 19 eventos... (copie e mude os IDs)

echo "Eventos enviados!"