fmt:
	go fmt ./...

vet:
	go vet ./...

build: fmt vet
	go build -o $(EXECUTABLE)

.PHONY: fmt vet build
