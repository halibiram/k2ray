.PHONY: test test-coverage

# Standard test runner
test:
	@echo "Running Go tests..."
	go test ./...

# Test runner with coverage report generation
test-coverage:
	@echo "Running Go tests with coverage..."
	go test -cover -coverprofile=coverage.out ./...

# Linter that outputs to a file for SonarQube
lint-report:
	@echo "Running golangci-lint and creating report..."
	@if ! command -v golangci-lint &> /dev/null; then \
		echo "golangci-lint not found, installing..."; \
		go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.59.1; \
	fi
	golangci-lint run --out-format json --issues-exit-code 0 ./... > golangci-report.json

# Add other build/run targets as needed
# e.g., build, clean, run
build:
	@echo "Building Go application..."
	go build -o k2ray-server ./cmd/k2ray

clean:
	@echo "Cleaning up build artifacts..."
	rm -f k2ray-server
	rm -f coverage.out

# ==============================================================================
# 📖 DOCUMENTATION
# ==============================================================================

swag:
	@echo ">> Generating Swagger/OpenAPI documentation..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@$(go env GOPATH)/bin/swag init -g cmd/k2ray/main.go