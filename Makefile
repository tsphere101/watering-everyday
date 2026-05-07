BINARY=watering-bot
GOOS=linux
GOARCH=amd64

.PHONY: build build-local test vet clean broadcast

build: vet test
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="-s -w" -o $(BINARY) ./src

build-local:
	go build -o $(BINARY) ./src
	go build -o broadcast ./cmd/broadcast

vet:
	go vet ./...

test:
	go test -count=1 ./...

broadcast:
	go build -o broadcast ./cmd/broadcast

clean:
	rm -f $(BINARY) broadcast

