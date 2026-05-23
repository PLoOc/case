#!/bin/bash

set -e

echo "Aguardando LocalStack ficar pronto..."
sleep 5

echo "Criando filas SQS..."

# Criar DLQ primeiro
awslocal sqs create-queue --queue-name raw-events-dlq 2>/dev/null || true
awslocal sqs create-queue --queue-name processed-events-dlq 2>/dev/null || true

# Criar filas principais com DLQ
awslocal sqs create-queue --queue-name raw-events \
  --attributes '{"RedrivePolicy":"{\"deadLetterTargetArn\":\"arn:aws:sqs:us-east-1:000000000000:raw-events-dlq\",\"maxReceiveCount\":\"3\"}"}'

awslocal sqs create-queue --queue-name processed-events \
  --attributes '{"RedrivePolicy":"{\"deadLetterTargetArn\":\"arn:aws:sqs:us-east-1:000000000000:processed-events-dlq\",\"maxReceiveCount\":\"3\"}"}'

echo "Criando tabelas DynamoDB..."

# Tabela events
awslocal dynamodb create-table \
  --table-name events \
  --attribute-definitions AttributeName=event_id,AttributeType=S \
  --key-schema AttributeName=event_id,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST 2>/dev/null || true

# Tabela developer_summary
awslocal dynamodb create-table \
  --table-name developer_summary \
  --attribute-definitions AttributeName=developer_id,AttributeType=S \
  --key-schema AttributeName=developer_id,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST 2>/dev/null || true

echo "Setup concluído!"