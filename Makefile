BINARY=watering-bot
GOOS=linux
GOARCH=amd64

.PHONY: build build-local test vet clean

build: vet test
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -ldflags="-s -w" -o $(BINARY) ./src

build-local:
	go build -o $(BINARY) ./src

vet:
	go vet ./...

test:
	go test -count=1 ./...

clean:
	rm -f $(BINARY)

