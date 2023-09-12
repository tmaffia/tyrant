.DEFAULT_GOAL := build-local

BINARY_NAME="discord-bots"


clean:
	go clean
	rm -f bin/*

build-local: clean
	GOARCH=amd64 GOOS=darwin go build -o bin/${BINARY_NAME}-darwin ./...
	GOARCH=amd64 GOOS=linux go build -o bin/${BINARY_NAME}-linux ./...
	GOARCH=amd64 GOOS=windows go build -o bin/${BINARY_NAME}.exe ./...

run: build-local
	bin/${BINARY_NAME}-linux

dep:
	go mod download