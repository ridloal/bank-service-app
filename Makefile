.PHONY: build run deploy clean migrate-up migrate-down

# Build the application
build:
	go build -o bin/main cmd/api/main.go

# Run the application locally
run:
	go run cmd/api/main.go

# Build and run docker containers
deploy:
	docker-compose up --build -d

# Stop and remove containers
clean:
	docker-compose down
	docker-compose rm -f

# Database migrations
migrate-up:
	docker-compose exec app migrate -path=/app/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@postgres:5432/${DB_NAME}?sslmode=disable" up

migrate-down:
	docker-compose exec app migrate -path=/app/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@postgres:5432/${DB_NAME}?sslmode=disable" down

# Show logs
logs:
	docker-compose logs -f

# Check status
status:
	docker-compose ps

# Restart services
restart:
	docker-compose restart