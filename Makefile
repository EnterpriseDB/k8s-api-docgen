.PHONY: lint all clean fmt k8s-api-docgen

all: k8s-api-docgen lint

k8s-api-docgen:
	go build -o bin/k8s-api-docgen cmd/k8s-api-docgen/main.go

lint:
	golangci-lint run

clean:
	rm -rf bin

fmt:
	go fmt ./...
