# Forge OS Development Makefile

.PHONY: help build test test-unit test-integration test-e2e test-coverage test-watch test-verbose benchmark metrics clean setup shell dev up down logs

# Default target
help: ## Show this help message
	@echo "Forge OS Development Commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

# Development environment
dev: ## Start development environment in Docker
	docker-compose up -d

up: dev ## Alias for dev

down: ## Stop development environment
	docker-compose down

shell: ## Open shell in development container
	docker-compose exec forge-dev bash

logs: ## Show development container logs
	docker-compose logs -f forge-dev

# Building
build: ## Build the Forge CLI binary
	go build -o bin/forge ./cmd/forge

build-linux: ## Build for Linux
	GOOS=linux GOARCH=amd64 go build -o bin/forge-linux ./cmd/forge

build-darwin: ## Build for macOS
	GOOS=darwin GOARCH=amd64 go build -o bin/forge-darwin ./cmd/forge

build-windows: ## Build for Windows
	GOOS=windows GOARCH=amd64 go build -o bin/forge-windows.exe ./cmd/forge

# Testing
test: test-unit test-integration ## Run all tests

test-unit: ## Run unit tests only
	gotestsum --format testname -- -race -v ./...

test-integration: ## Run integration tests
	gotestsum --format testname -- -race -v -tags=integration ./...

test-e2e: ## Run end-to-end tests
	gotestsum --format testname -- -race -v -tags=e2e ./test/e2e/

test-coverage: ## Generate coverage report
	gotestsum --format testname -- -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-watch: ## Run tests on file changes (requires entr)
	find . -name "*.go" -not -path "./vendor/*" | entr -r make test-unit

test-verbose: ## Run tests with verbose output
	gotestsum --format standard-verbose -- -race -v ./...

# Benchmarking
benchmark: ## Run performance benchmarks
	go test -bench=. -benchmem ./... | tee benchmarks/latest.txt

benchmark-compare: ## Compare benchmarks with previous run
	@echo "Comparing benchmarks..."
	@go test -bench=. -benchmem ./... > benchmarks/current.txt
	@benchstat benchmarks/previous.txt benchmarks/current.txt || echo "No previous benchmark to compare"

# Metrics
metrics: ## Generate performance metrics report
	@echo "Generating metrics report..."
	@mkdir -p metrics
	@echo "# Forge OS Metrics Report" > metrics/report.md
	@echo "Generated: $$(date)" >> metrics/report.md
	@echo "" >> metrics/report.md
	@echo "## Build Time" >> metrics/report.md
	@time -p make build 2>&1 | grep real | awk '{print "Build time:", $$2, "seconds"}' >> metrics/report.md
	@echo "" >> metrics/report.md
	@echo "## Test Execution Time" >> metrics/report.md
	@time -p make test-unit 2>&1 | grep real | awk '{print "Unit tests:", $$2, "seconds"}' >> metrics/report.md
	@echo "" >> metrics/report.md
	@echo "## Binary Size" >> metrics/report.md
	@make build > /dev/null 2>&1
	@ls -lh bin/forge | awk '{print "Binary size:", $$5}' >> metrics/report.md
	@echo "" >> metrics/report.md
	@echo "## Code Statistics" >> metrics/report.md
	@find . -name "*.go" -not -path "./vendor/*" | wc -l | awk '{print "Go files:", $$1}' >> metrics/report.md
	@grep -r "^func" --include="*.go" . --exclude-dir=vendor | wc -l | awk '{print "Functions:", $$1}' >> metrics/report.md
	@echo "Metrics report generated: metrics/report.md"

# Setup and cleanup
setup: ## Set up development environment
	go mod download
	go mod tidy
	mkdir -p bin output metrics benchmarks
	@echo "Development environment set up"

clean: ## Clean build artifacts and caches
	rm -rf bin/ output/ build/ dist/ *.img *.qcow2 *.raw
	rm -f coverage.out coverage.html
	docker-compose down -v 2>/dev/null || true
	docker system prune -f 2>/dev/null || true

clean-all: clean ## Clean everything including caches
	rm -rf dl/ .ccache/ .cache/ .forge/
	rm -rf metrics/ benchmarks/
	go clean -cache
	go clean -modcache
	docker-compose down -v
	docker system prune -f
	docker volume prune -f

# Docker utilities
docker-build: ## Rebuild Docker development environment
	docker-compose build --no-cache

docker-clean: ## Clean Docker containers and volumes
	docker-compose down -v
	docker system prune -f
	docker volume prune -f

# Git utilities
git-status: ## Show git status
	@echo "=== Git Status ==="
	git status --short
	@echo ""
	@echo "=== Recent Commits ==="
	git log --oneline -5

git-clean: ## Clean untracked files (DANGER: removes uncommitted work)
	git clean -fd
	git reset --hard HEAD

# CI/CD simulation
ci: ## Simulate CI pipeline
	@echo "=== Running CI Pipeline ==="
	make setup
	make test-coverage
	make build
	make benchmark
	@echo "CI Pipeline completed successfully"