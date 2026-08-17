BINARY_NAME := slashviberepo
GO := go

.PHONY: build test lint clean

build:
	CGO_ENABLED=0 $(GO) build -a -installsuffix cgo -o $(BINARY_NAME) .

test:
	$(GO) test -v -race -coverprofile=coverage.out ./...

lint:
	golangci-lint run

clean:
	rm -f $(BINARY_NAME) coverage.out
