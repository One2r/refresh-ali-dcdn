# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
BINARY_NAME=refresh-ali-dcdn
BINARY_DARWIN=$(BINARY_NAME).darwin
BINARY_UNIX=$(BINARY_NAME).unix
BINARY_ARM=$(BINARY_NAME).arm64
BINARY_WIN=$(BINARY_NAME).exe

build:
		$(GOBUILD) -o $(BINARY_NAME) -v

clean:
		$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_DARWIN)
		rm -f $(BINARY_UNIX)
		rm -f $(BINARY_ARM)

# Cross compilation
build-darwin:
		CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_DARWIN) -v
build-unix:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
build-arm64:
		CGO_ENABLED=0 GOOS=linux GOARCH=arm64 $(GOBUILD) -o $(BINARY_ARM) -v
build-win:
		CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_WIN) -v
