create_mocks:
	mockery --all --recursive --output ./mocks

linter:
	golangci-lint run internal/... --timeout 5m