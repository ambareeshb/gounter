APP_NAME = gounter

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	docker-compose -f ./infra/docker-compose.yml build

# Run the application
run: build
	@echo "Running $(APP_NAME)..."
	docker-compose -f ./infra/docker-compose.yml up gounter

# Test the application
test:
	@echo "Running tests..."
	go test ./...

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	docker-compose -f ./infra/docker-compose.yml down

# Create db migration file
create-migration:
	migrate create  -dir infra/db/migrations/ -ext sql ${name}
# Phony targets
.PHONY: build run test clean