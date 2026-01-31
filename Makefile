.PHONY: test test-cover test-race lint build run clean coverage-html coverage-func coverage-check

# Test targets
test:
	go test -v ./...

test-cover:
	go test -v -coverprofile=coverage.out -covermode=atomic ./...

test-race:
	go test -v -race ./...

# Coverage reporting
coverage-html:
	go test -coverprofile=coverage.out -covermode=atomic ./... && \
	go tool cover -html=coverage.out -o coverage.html

coverage-func:
	go test -coverprofile=coverage.out -covermode=atomic ./... && \
	go tool cover -func=coverage.out

coverage-check:
	@go test -coverprofile=coverage.out -covermode=atomic ./... && \
	TOTAL=$$(go tool cover -func=coverage.out | tail -1 | awk '{print $$3}' | sed 's/%//'); \
	echo "Total coverage: $$TOTAL%"; \
	MIN_COVERAGE=60; \
	if [ $$(echo "$$TOTAL < $$MIN_COVERAGE" | bc) -eq 1 ]; then \
		echo "❌ Coverage ($$TOTAL%) is below threshold ($$MIN_COVERAGE%)"; \
		exit 1; \
	else \
		echo "✅ Coverage ($$TOTAL%) meets threshold ($$MIN_COVERAGE%)"; \
	fi

coverage-packages:
	@go test -coverprofile=coverage.out -covermode=atomic ./... && \
	echo "=== Per-package coverage ===" && \
	go tool cover -func=coverage.out | grep "^github.com" | awk '{
		pkg = $$1;
		gsub(/\/github.com\/atumaikin\/nexflow\//, "", pkg);
		cov = $$3;
		print pkg ": " cov
	}' | sort

# Linting
lint:
	golangci-lint run --timeout=5m

vet:
	go vet ./...

# Build targets
build:
	go build -v ./...

run:
	go run cmd/server/main.go

# Cleanup
clean:
	@echo "Cleaning up..."
	@find . -name "coverage.out" -delete
	@find . -name "coverage.html" -delete
	@find . -name "coverage-*.txt" -delete
	@find . -name "*.db" -delete
	@find . -name "*.sqlite" -delete
	@find . -name "*.sqlite3" -delete
	@echo "Cleanup complete!"

# Dependencies
deps:
	go mod download
	go mod verify
	go mod tidy

# Format
fmt:
	go fmt ./...
	gofmt -s -w .

# All checks
ci: test-cover lint vet coverage-check

# Development helpers
dev:
	go run cmd/server/main.go
