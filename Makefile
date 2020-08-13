run:
	docker build -t avbru/image-previewer-server .
	docker-compose -p previewer -f ./deployments/docker-compose.yaml up previewer
test:
	go test -race -count 2 ./...
run_integration_test:
	./scripts/run-integration-tests.sh
lint:
	gofumpt -s -w .
	golangci-lint run ./...

.PHONY: all test clean