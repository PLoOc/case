#!/bin/bash

set -e

echo "Criando filas SQS..."

# DLQ para raw-events
awslocal sqs create-queue --queue-name raw-events-dlq

# Fila raw-events com redirecionamento para DLQ
awslocal sqs create-queue --queue-name raw-events \
  --attributes '{"RedrivePolicy": "{\"deadLetterTargetArn\":\"arn:aws:sqs:us-east-1:000000000000:raw-events-dlq\",\"maxReceiveCount\":\"3\"}"}'

# DLQ para o processamento
awslocal sqs create-queue --queue-name processed-events-dlq

# Fila que processa e redireciona pro DLQ
awslocal sqs create-queue --queue-name processed-events \
  --attributes '{"RedrivePolicy": "{\"deadLetterTargetArn\":\"arn:aws:sqs:us-east-1:000000000000:processed-events-dlq\",\"maxReceiveCount\":\"3\"}"}'

echo "Criando tabelas DynamoDB..."

# Tabela dos eventos
awslocal dynamodb create-table \
  --table-name events \
  --attribute-definitions AttributeName=event_id,AttributeType=S \
  --key-schema AttributeName=event_id,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST

# Tabela por dev
awslocal dynamodb create-table \
  --table-name developer_summary \
  --attribute-definitions AttributeName=developer_id,AttributeType=S \
  --key-schema AttributeName=developer_id,KeyType=HASH \
  --billing-mode PAY_PER_REQUEST

echo "Setup concluído!"