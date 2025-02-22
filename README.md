# Bank Service

A REST API service for basic banking operations implemented using Go, Echo framework, and PostgreSQL.

## Features

- Clean Architecture implementation
- REST API endpoints for banking operations
- Structured logging with Zap
- PostgreSQL database with migrations
- Docker and Docker Compose setup
- Configuration management
- Comprehensive error handling

## Tech Stack

- Go 1.21
- Echo Framework
- PostgreSQL
- Docker & Docker Compose
- Zap Logger
- Golang Migrate
- Viper (Configuration)

## Project Structure

```
bank-service/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── domain/
│   ├── repository/
│   ├── usecase/
│   └── delivery/
│       └── http/
├── pkg/
│   ├── config/
│   └── logger/
├── migrations/
├── Dockerfile
├── docker-compose.yml
└── Makefile
```

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- Make

## Getting Started

1. Clone the repository:
```bash
git clone https://github.com/ridloal/bank-service-app
cd bank-service-app
```

2. Copy environment file:
```bash
cp .env.example .env
```

3. Build and run with Docker:
```bash
make deploy
```

4. Run database migrations:
```bash
make migrate-up
```

5. Run tests:
```bash
make test
```

## Development

1. Run locally without Docker:
```bash
go mod download
go run cmd/api/main.go
```

2. Build binary:
```bash
make build
```

3. Run tests:
```bash
go test ./...
```

## Configuration

### Environment Variables

- `APP_DB_HOST`: Database host
- `APP_DB_PORT`: Database port
- `APP_DB_USER`: Database user
- `APP_DB_PASSWORD`: Database password
- `APP_DB_NAME`: Database name
- `APP_DB_SSLMODE`: SSL mode for database connection

### Command Line Arguments

- `--server.host`: Server host (default: "0.0.0.0")
- `--server.port`: Server port (default: "8080")
- `--server.read_timeout`: Server read timeout (default: 15s)
- `--server.write_timeout`: Server write timeout (default: 15s)

## Deployment

### With Docker

1. Build and start services:
```bash
make deploy
```

2. Check service status:
```bash
make status
```

3. View logs:
```bash
make logs
```

4. Stop services:
```bash
make clean
```

### Manual Deployment

1. Build the binary:
```bash
make build
```

2. Configure environment variables.

3. Run the binary:
```bash
./bin/main
```

## Database Migrations

- Run migrations up:
```bash
make migrate-up
```

- Rollback migrations:
```bash
make migrate-down
```

## API Documentation

- See [API Documentation](docs/API.md) for detailed endpoint information.
- Postman Collection : [Postman Collection](docs/BankPostmanCollection.postman_collection.json)
- Postman Collection Web : https://documenter.getpostman.com/view/15911180/2sAYX3r3mi

## Logging

The application uses structured logging with the following levels:
- INFO: Normal operations
- WARN: Warning conditions
- ERROR: Error conditions
- DEBUG: Debug information

Logs are output in JSON format for easy parsing.