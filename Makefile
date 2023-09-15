.DEFAULT_GOAL := build

BINARY_NAME="discord-bots"

ifndef ${TAG}
TAG = latest
endif


clean:
	go clean
	rm -f bin/*

build: clean
	GOARCH=amd64 GOOS=darwin go build -o bin/${BINARY_NAME}-darwin ./...
	GOARCH=amd64 GOOS=linux go build -o bin/${BINARY_NAME}-linux ./...
	GOARCH=amd64 GOOS=windows go build -o bin/${BINARY_NAME}.exe ./...

buildc: clean
	docker build -f Dockerfile -t szerated/${BINARY_NAME}:${TAG} .

run: build
	bin/${BINARY_NAME}-linux

runc: buildc
		docker run -it szerated/${BINARY_NAME}:${TAG}

dep:
	go mod download