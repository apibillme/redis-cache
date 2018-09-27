# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
BINARY_NAME=redis-cache

test: test_internal build
build: 
	$(GOBUILD) -o $(BINARY_NAME) -v
test_internal: 
	$(GOTEST) -v ./...
    