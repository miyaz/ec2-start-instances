# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
DEPCMD=dep
DEPENSURE=$(DEPCMD) ensure
BINARY_NAME=ec2start
BINARY_MAC=$(BINARY_NAME)_macos
BINARY_ZIP=$(BINARY_NAME).zip

all: deps test build
test: 
	$(GOTEST) -v ./...
clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_MAC)
	rm -f $(BINARY_ZIP)
run:
	./$(BINARY_NAME)
deps:
	$(DEPENSURE)
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) -v main.go
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_MAC) -v main.go
	zip $(BINARY_ZIP) ./$(BINARY_NAME)

