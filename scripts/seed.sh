#!/bin/bash

echo "ENVIANDO EVENTOS DE TESTE..."

# ===== EVENTOS VÁLIDOS =====

# dev-123 - commits
docker exec case-localstack awslocal sqs send-message --queue-url http://localhost:4566/000000000000/raw-events --message-body '{"event_id":"550e8400-e29b-41d4-a716-446655440001","developer_id":"dev-123","metric_type":"commits","value":5,"repository":"org/repo-1","timestamp":"2026-05-21T10:00:00Z"}'
echo "✅ Evento 1 enviado"

docker exec case-localstack awslocal sqs send-message --queue-url http://localhost:4566/000000000000/raw-events --message-body '{"event_id":"550e8400-e29b-41d4-a716-446655440002","developer_id":"dev-123","metric_type":"commits","value":8,"repository":"org/repo-2","timestamp":"2026-05-21T11:00:00Z"}'
echo "✅ Evento 2 enviado"

docker exec case-localstack awslocal sqs send-message --queue-url http://localhost:4566/000000000000/raw-events --message-body '{"event_id":"550e8400-e29b-41d4-a716-446655440003","developer_id":"dev-123","metric_type":"pull_requests","value":3,"repository":"org/repo-1","timestamp":"2026-05-21T12:00:00Z"}'
echo "✅ Evento 3 enviado"

docker exec case-localstack awslocal sqs send-message --queue-url http://localhost:4566/000000000000/raw-events --message-body '{"event_id":"550e8400-e29b-41d4-a716-446655440004","developer_id":"dev-123","metric_type":"review_time_minutes","value":45,"repository":"org/repo-3","timestamp":"2026-05-21T13:00:00Z"}'
echo "✅ Evento 4 enviado"

docker exec case-localstack awslocal sqs send-message --queue-url http://localhost:4566/000000000000/raw-events --message-body '{"event_id":"550e8400-e29b-41d4-a716-446655440005","developer_id":"dev-123","metric_type":"review_time_minutes","value":90,"repository":"org/repo-2","timestamp":"2026-05-21T14:00:00Z"}'
echo "✅ Evento 5 enviado"

# dev-456 - commits e pull_requests
docker exec case-localstack awslocal sqs send-message --queue-url http://localhost:4566/000000000000/raw-events --message-body '{"event_id":"550e8400-e29b-41d4-a716-446655440006","developer_id":"dev-456","metric_type":"commits","value":12,"repository":"org/repo-1","timestamp":"2026-05-21T10:30:00Z"}'
echo "✅ Evento 6 enviado"

docker exec case-localstack awslocal sqs send-message --queue-url http://localhost:4566/000000000000/raw-events --message-body '{"event_id":"550e8400-e29b-41d4-a716-446655440007","developer_id":"dev-456","metric_type":"commits","value":7,"repository":"org/repo-4","timestamp":"2026-05-21T11:30:00Z"}'
echo "✅ Evento 7 enviado"

docker exec case-localstack awslocal sqs send-message --queue-url http://localhost:4566/000000000000/raw-events --message-body '{"event_id":"550e8400-e29b-41d4-a716-446655440008","developer_id":"dev-456","metric_type":"pull_requests","value":5,"repository":"org/repo-1","timestamp":"2026-05-21T12:30:00Z"}'
echo "✅ Evento 8 enviado"

docker exec case-localstack awslocal sqs send-message --queue-url http://localhost:4566/000000000000/raw-events --message-body '{"event_id":"550e8400-e29b-41d4-a716-446655440009","developer_id":"dev-456","metric_type":"review_time_minutes","value":120,"repository":"org/repo-2","timestamp":"2026-05-21T13:30:00Z"}'
echo "✅ Evento 9 enviado"

# dev-789
docker exec case-localstack awslocal sqs send-message --queue-url http://localhost:4566/000000000000/raw-events --message-body '{"event_id":"550e8400-e29b-41d4-a716-446655440010","developer_id":"dev-789","metric_type":"commits","value":20,"repository":"org/repo-5","timestamp":"2026-05-21T09:00:00Z"}'
echo "✅ Evento 10 enviado"

docker exec case-localstack awslocal sqs send-message --queue-url http://localhost:4566/000000000000/raw-events --message-body '{"event_id":"550e8400-e29b-41d4-a716-446655440011","developer_id":"dev-789","metric_type":"pull_requests","value":10,"repository":"org/repo-5","timestamp":"2026-05-21T10:00:00Z"}'
echo "✅ Evento 11 enviado"

docker exec case-localstack awslocal sqs send-message --queue-url http://localhost:4566/000000000000/raw-events --message-body '{"event_id":"550e8400-e29b-41d4-a716-446655440012","developer_id":"dev-789","metric_type":"review_time_minutes","value":30,"repository":"org/repo-3","timestamp":"2026-05-21T11:00:00Z"}'
echo "✅ Evento 12 enviado"

docker exec case-localstack awslocal sqs send-message --queue-url http://localhost:4566/000000000000/raw-events --message-body '{"event_id":"550e8400-e29b-41d4-a716-446655440013","developer_id":"dev-789","metric_type":"commits","value":4,"repository":"org/repo-1","timestamp":"2026-05-21T15:00:00Z"}'
echo "✅ Evento 13 enviado"

docker exec case-localstack awslocal sqs send-message --queue-url http://localhost:4566/000000000000/raw-events --message-body '{"event_id":"550e8400-e29b-41d4-a716-446655440014","developer_id":"dev-789","metric_type":"pull_requests","value":2,"repository":"org/repo-4","timestamp":"2026-05-21T16:00:00Z"}'
echo "✅ Evento 14 enviado"

# dev-999 - review_time_minutes no limite (1440 = 24h)
docker exec case-localstack awslocal sqs send-message --queue-url http://localhost:4566/000000000000/raw-events --message-body '{"event_id":"550e8400-e29b-41d4-a716-446655440015","developer_id":"dev-999","metric_type":"review_time_minutes","value":1440,"repository":"org/repo-6","timestamp":"2026-05-21T08:00:00Z"}'
echo "✅ Evento 15 enviado"

docker exec case-localstack awslocal sqs send-message --queue-url http://localhost:4566/000000000000/raw-events --message-body '{"event_id":"550e8400-e29b-41d4-a716-446655440016","developer_id":"dev-999","metric_type":"commits","value":1,"repository":"org/repo-6","timestamp":"2026-05-21T09:00:00Z"}'
echo "✅ Evento 16 enviado"

# ===== DUPLICADOS (idempotência) =====

docker exec case-localstack awslocal sqs send-message --queue-url http://localhost:4566/000000000000/raw-events --message-body '{"event_id":"550e8400-e29b-41d4-a716-446655440001","developer_id":"dev-123","metric_type":"commits","value":5,"repository":"org/repo-1","timestamp":"2026-05-21T10:00:00Z"}'
echo "⚠️  Evento 17 enviado (DUPLICADO do evento 1 - deve ser ignorado)"

docker exec case-localstack awslocal sqs send-message --queue-url http://localhost:4566/000000000000/raw-events --message-body '{"event_id":"550e8400-e29b-41d4-a716-446655440006","developer_id":"dev-456","metric_type":"commits","value":12,"repository":"org/repo-1","timestamp":"2026-05-21T10:30:00Z"}'
echo "⚠️  Evento 18 enviado (DUPLICADO do evento 6 - deve ser ignorado)"

# ===== INVÁLIDOS (devem ir para DLQ) =====

docker exec case-localstack awslocal sqs send-message --queue-url http://localhost:4566/000000000000/raw-events --message-body '{"event_id":"nao-e-um-uuid","developer_id":"dev-123","metric_type":"commits","value":5,"repository":"org/repo-1","timestamp":"2026-05-21T10:00:00Z"}'
echo "❌ Evento 19 enviado (INVÁLIDO - event_id não é UUID)"

docker exec case-localstack awslocal sqs send-message --queue-url http://localhost:4566/000000000000/raw-events --message-body '{"event_id":"550e8400-e29b-41d4-a716-446655440020","developer_id":"","metric_type":"commits","value":5,"repository":"org/repo-1","timestamp":"2026-05-21T10:00:00Z"}'
echo "❌ Evento 20 enviado (INVÁLIDO - developer_id vazio)"

docker exec case-localstack awslocal sqs send-message --queue-url http://localhost:4566/000000000000/raw-events --message-body '{"event_id":"550e8400-e29b-41d4-a716-446655440021","developer_id":"dev-123","metric_type":"tipo_errado","value":5,"repository":"org/repo-1","timestamp":"2026-05-21T10:00:00Z"}'
echo "❌ Evento 21 enviado (INVÁLIDO - metric_type inválido)"

docker exec case-localstack awslocal sqs send-message --queue-url http://localhost:4566/000000000000/raw-events --message-body '{"event_id":"550e8400-e29b-41d4-a716-446655440022","developer_id":"dev-123","metric_type":"review_time_minutes","value":9999,"repository":"org/repo-1","timestamp":"2026-05-21T10:00:00Z"}'
echo "❌ Evento 22 enviado (INVÁLIDO - review_time_minutes acima de 1440)"

docker exec case-localstack awslocal sqs send-message --queue-url http://localhost:4566/000000000000/raw-events --message-body '{"event_id":"550e8400-e29b-41d4-a716-446655440023","developer_id":"dev-123","metric_type":"commits","value":5,"repository":"org/repo-1","timestamp":"2099-01-01T00:00:00Z"}'
echo "❌ Evento 23 enviado (INVÁLIDO - timestamp futuro)"

echo ""
echo "✅ Seed concluído! 16 válidos, 2 duplicados, 5 inválidos"
echo ""
echo "Aguarde alguns segundos e consulte:"
echo "  curl http://localhost:8080/metrics/dev-123/summary"
echo "  curl http://localhost:8080/metrics/dev-456/summary"
echo "  curl http://localhost:8080/metrics/dev-789/summary"
echo "  curl http://localhost:8080/metrics/dev-999/summary"