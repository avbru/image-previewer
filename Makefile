run:
	docker build -t avbru/image-previewer-server .
	docker-compose -p previewer -f ./deployments/docker-compose.yaml up previewer
test:
	go test -race -count 10 ./...
run_integration_test:
	docker build -f integration-test.Dockerfile -t avbru/integration-tests .
	docker build -t avbru/fileserver ./fileserver
	docker build -t avbru/image-previewer-server .
	./scripts/run-integration-tests.sh
lint:
	golangci-lint run ./...

.PHONY: all test clean