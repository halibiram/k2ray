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
# 🏠 KEENETIC ROUTER BUILDS
# ==============================================================================

# Build for Keenetic Extra DSL KN2112 (MIPS little-endian)
build-keenetic:
	@echo "Building K2Ray for Keenetic Extra DSL KN2112..."
	@chmod +x deployments/entware/scripts/build_for_keenetic.sh
	@./deployments/entware/scripts/build_for_keenetic.sh

# Build for all supported MIPS architectures
build-keenetic-all:
	@echo "Building K2Ray for all supported Keenetic architectures..."
	@chmod +x deployments/entware/scripts/build_for_keenetic.sh
	@./deployments/entware/scripts/build_for_keenetic.sh --all-archs

# Cross-compile for MIPS manually
build-mips:
	@echo "Cross-compiling for MIPS (little-endian)..."
	@mkdir -p build/mips
	@GOOS=linux GOARCH=mipsle CGO_ENABLED=0 go build \
		-tags "netgo,osusergo" \
		-ldflags "-s -w -extldflags '-static'" \
		-o build/mips/k2ray ./cmd/k2ray

# Cross-compile for big-endian MIPS
build-mips-be:
	@echo "Cross-compiling for MIPS (big-endian)..."
	@mkdir -p build/mips-be
	@GOOS=linux GOARCH=mips CGO_ENABLED=0 go build \
		-tags "netgo,osusergo" \
		-ldflags "-s -w -extldflags '-static'" \
		-o build/mips-be/k2ray ./cmd/k2ray

# Install Entware package locally (for testing)
install-entware-local:
	@echo "Installing K2Ray for Entware (local testing)..."
	@if [ ! -f build/mips/k2ray ]; then \
		echo "Building MIPS binary first..."; \
		make build-mips; \
	fi
	@chmod +x deployments/entware/scripts/install_entware.sh
	@sudo ./deployments/entware/scripts/install_entware.sh

# Create Keenetic installation package
package-keenetic:
	@echo "Creating Keenetic installation package..."
	@make build-keenetic
	@echo "Package created in build/keenetic/ directory"

# ==============================================================================
# 📖 DOCUMENTATION
# ==============================================================================

swag:
	@echo ">> Generating Swagger/OpenAPI documentation..."
	@go install github.com/swaggo/swag/cmd/swag@latest
	@$(go env GOPATH)/bin/swag init -g cmd/k2ray/main.go