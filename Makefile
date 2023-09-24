.DEFAULT_GOAL := build

BINARY_NAME=tyrant
DOCKERHUB_USER=szerated
CLEAN_CONTAINERS=docker ps -a | awk '{ print $$1,$$2 }' | grep ${BINARY_NAME} | awk '{print $$1 }' | xargs -I {} docker rm {}


ifndef ${TAG}
	TAG = latest
endif


clean:
	go clean
	rm -f bin/*

cleanc: clean
	${CLEAN_CONTAINERS}
	docker rmi -f ${DOCKERHUB_USER}/${BINARY_NAME}:${TAG}

build: clean
	GOARCH=amd64 GOOS=darwin go build -o bin/${BINARY_NAME}-darwin ./...
	GOARCH=amd64 GOOS=linux go build -o bin/${BINARY_NAME}-linux ./...
	GOARCH=amd64 GOOS=windows go build -o bin/${BINARY_NAME}.exe ./...

buildc: cleanc
	docker build -f Dockerfile -t ${DOCKERHUB_USER}/${BINARY_NAME}:${TAG} .

run: build
	bin/${BINARY_NAME}-linux

runc: buildc
	docker run -it ${DOCKERHUB_USER}/${BINARY_NAME}:${TAG}

dep:
	go mod download