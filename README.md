# Developer Metrics Pipeline

Pipeline de métricas de produtividade de desenvolvedores construído com dois serviços Go que se comunicam via filas SQS, usando LocalStack para simular a infraestrutura AWS localmente.

```
[raw-events] → Processor → [processed-events] → Aggregator → DynamoDB → API REST
```

## Pré-requisitos

- Docker
- Docker Compose


## Subindo o projeto

```bash
docker-compose up --build
```

Isso sobe quatro containers em ordem:

1. **localstack** — AWS local (SQS + DynamoDB)
2. **setup** — cria as filas e tabelas automaticamente
3. **processor** — começa a consumir a fila `raw-events`
4. **aggregator** — começa a consumir a fila `processed-events` e sobe a API na porta `8080`

## Enviando eventos

Script de seed, com cenários para popular com dados de teste:

```bash
bash scripts/seed.sh
```

O script envia 23 mensagens: 16 válidas (4 developers diferentes), 2 duplicadas e 5 inválidas — cobrindo os cenários de validação, idempotência e DLQ.

Ou manualmente via container do LocalStack:

```bash
docker exec case-localstack awslocal sqs send-message \
  --queue-url http://localhost:4566/000000000000/raw-events \
  --message-body '{
    "event_id": "550e8400-e29b-41d4-a716-446655440001",
    "developer_id": "dev-123",
    "metric_type": "commits",
    "value": 5,
    "repository": "org/repo-1",
    "timestamp": "2026-05-21T10:00:00Z"
  }'
```

## API

### Health check
```
GET /health
```

### Eventos de um desenvolvedor
```
GET /metrics/{developer_id}
```

### Resumo agregado
```
GET /metrics/{developer_id}/summary
```

Exemplo de resposta:
```json
{
  "developer_id": "dev-123",
  "total_commits": 13,
  "total_pull_requests": 3,
  "total_review_time": 135,
  "avg_review_time_minutes": 67.5,
  "events_processed": 5,
  "last_activity": "2026-05-21T14:00:00Z"
}
```

## Estrutura do projeto

```
.
├── services/
│   ├── processor/
│   │   ├── cmd/main.go
│   │   ├── internal/
│   │   │   ├── domain/       # validação e regras de negócio
│   │   │   └── infra/
│   │   │       ├── queue/    # adapter SQS
│   │   │       └── worker/   # worker pool concorrente
│   │   └── Dockerfile
│   └── aggregator/
│       ├── cmd/main.go
│       ├── internal/
│       │   ├── domain/       # entidades e regras de agregação
│       │   ├── usecase/      # orquestração
│       │   └── infra/
│       │       ├── api/      # handlers HTTP
│       │       ├── queue/    # adapter SQS
│       │       └── repository/ # adapter DynamoDB
│       └── Dockerfile
├── infra/
│   └── localstack/
│       └── init-aws.sh
├── scripts/
│   └── seed.sh
└── docker-compose.yml
```

## Variáveis de ambiente

### Processor

| Variável | Padrão | Descrição |
|---|---|---|
| `AWS_ENDPOINT_URL` | `http://localhost:4566` | Endpoint da AWS/LocalStack |
| `AWS_REGION` | `us-east-1` | Região |
| `RAW_EVENTS_QUEUE` | `raw-events` | Fila de entrada |
| `PROCESSED_EVENTS_QUEUE` | `processed-events` | Fila de saída |
| `WORKER_COUNT` | `3` | Número de workers concorrentes |

### Aggregator

| Variável | Padrão | Descrição |
|---|---|---|
| `AWS_ENDPOINT_URL` | `http://localhost:4566` | Endpoint da AWS/LocalStack |
| `AWS_REGION` | `us-east-1` | Região |
| `PROCESSED_EVENTS_QUEUE` | `processed-events` | Fila consumida |
| `PORT` | `8080` | Porta da API |

## Validações do Processor

Um evento é rejeitado (vai para `raw-events-dlq` após 3 tentativas) se:

- `event_id` estiver vazio ou não for um UUID válido
- `developer_id` estiver vazio
- `metric_type` não for `commits`, `pull_requests` ou `review_time_minutes`
- `value` for negativo
- `value` for maior que `1440` quando `metric_type` for `review_time_minutes`
- `timestamp` for uma data futura

## Verificando a DLQ

Após envio dos eventos inválidos, aguardar ~90 segundos (30s de visibility timeout × 3 tentativas) e verificar:

```bash
docker exec case-localstack awslocal sqs get-queue-attributes \
  --queue-url http://localhost:4566/000000000000/raw-events-dlq \
  --attribute-names ApproximateNumberOfMessages
```
