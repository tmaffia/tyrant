.DEFAULT_GOAL := build-local

BINARY_NAME="stop-bot"

ifndef ${TAG}
TAG = latest
endif

clean:
	go clean
	rm -f bin/*

build-local: clean
	GOARCH=amd64 GOOS=darwin go build -o bin/${BINARY_NAME}-darwin ./...
	GOARCH=amd64 GOOS=linux go build -o bin/${BINARY_NAME}-linux ./...
	GOARCH=amd64 GOOS=windows go build -o bin/${BINARY_NAME}.exe ./...

build-container: clean
	docker build -f build/Dockerfile -t szerated/${BINARY_NAME}:${TAG} .

run: build-local
	bin/${BINARY_NAME}-linux

run-container: build-container
	docker run --env-file .env -it szerated/${BINARY_NAME}:${TAG}

dep:
	go mod download