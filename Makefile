fmt:
	go fmt ./...

vet:
	go vet ./...

build: fmt vet
	go build -o $(EXECUTABLE) ./cmd/learn_db

.PHONY: fmt vet build
