APP_NAME = gounter
RICHGO = richgo

# Help target to display usage information
help:
	@echo "Makefile commands for $(APP_NAME):"
	@echo ""
	@echo "  build               Build the application"
	@echo "  run                 Run the application"
	@echo "  test                Run unit tests"
	@echo "  test-integration    Run integration tests"
	@echo "  install-deps        Install richgo if not available"
	@echo "  clean               Clean up build artifacts"
	@echo "  create-migration    Create a new database migration file"
	@echo ""

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
	@echo "Running unit tests..."
	@$(RICHGO) test ./internal/...
	@$(RICHGO) test ./api/...

# run integration test for the application, this need the application to be up and running
test-integration:
	@echo "Running integration tests..."
	@$(RICHGO) test ./test/integration/...

# Install richgo if not available
install-deps:
	@which richgo > /dev/null || { \
		go install github.com/kyoh86/richgo@latest; \
	}

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	docker-compose -f ./infra/docker-compose.yml down
	rm -rf ./infra/postgres-data

# Create db migration file
create-migration:
	migrate create  -dir infra/db/migrations/ -ext sql ${name}

# serve swagger documentation at 8080 port
serve-swagger:
	docker-compose -f ./docs/docker-compose.yml up -d  --build
	@open http://localhost:8080 || xdg-open http://localhost:8080 || start http://localhost:8080

# Phony targets
.PHONY: build run test clean