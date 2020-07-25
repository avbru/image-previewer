build:
	go build -o previewer ./...
run:
	@echo "  >  starting previewer server"
	@-./previewer 2>&1 & echo $$! > $(PID)
test:
	go test -race ./...
lint:
	golangci-lint run ./...