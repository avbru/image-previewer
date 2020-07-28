FROM golang:1.14-alpine AS builder

WORKDIR /app
COPY . .

RUN GOOS=linux GOARCH=amd64 go build -o /go/bin/previewer ./cmd/previewer/...
ENTRYPOINT ["/go/bin/previewer"]
EXPOSE 80
