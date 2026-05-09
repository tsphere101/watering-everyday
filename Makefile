BINARY=watering-bot
BINARY_MACOS=watering-bot-macos
GOOS=linux
GOARCH=amd64

.PHONY: build test vet clean broadcast

build: vet test
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="-s -w" -o $(BINARY) ./src
	go build -o $(BINARY_MACOS) ./src

vet:
	go vet ./...

test:
	go test -count=1 ./...

broadcast:
	go build -o broadcast ./cmd/broadcast

clean:
	rm -f $(BINARY) $(BINARY_MACOS) broadcast

