.PHONY: fmt lint check install-hooks

# Format all Go files
fmt:
	gofmt -w -s .

# Run linter
lint:
	golangci-lint run

# Format and lint check
check: fmt lint

# Install Git hooks
install-hooks:
	cp scripts/pre-commit.sh .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit

# Run tests
test:
	go test ./...