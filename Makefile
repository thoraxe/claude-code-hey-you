.PHONY: build build-all clean test

BINARY_NAME=claude-code-hey-you
LDFLAGS=-ldflags "-s -w"

# Build for current platform
build:
	go build $(LDFLAGS) -o $(BINARY_NAME)

# Build for all platforms
build-all: clean
	mkdir -p dist
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-amd64.exe
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-windows-arm64.exe
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-amd64
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-linux-arm64
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-amd64
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o dist/$(BINARY_NAME)-darwin-arm64

clean:
	rm -rf dist/
	rm -f $(BINARY_NAME) $(BINARY_NAME).exe

test:
	go test -v ./...

# Test with sample input
test-stop:
	echo '{"hook_event_name":"Stop","cwd":"/test/project"}' | ./$(BINARY_NAME) --topic test-topic

test-bash:
	echo '{"hook_event_name":"PreToolUse","tool_name":"Bash","tool_input":{"command":"rm -rf ./temp"}}' | ./$(BINARY_NAME) --topic test-topic

test-notification:
	echo '{"hook_event_name":"Notification","message":"Claude needs your permission"}' | ./$(BINARY_NAME) --topic test-topic
