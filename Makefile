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
	go build -gcflags "-N -l" -o bin/${BINARY_NAME} ./...

buildc: cleanc
	docker build -f Dockerfile -t ${DOCKERHUB_USER}/${BINARY_NAME}:${TAG} .

run: build
	bin/${BINARY_NAME}

runc: buildc
	docker run -it ${DOCKERHUB_USER}/${BINARY_NAME}:${TAG}

dep:
	go mod download