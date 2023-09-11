FROM golang:1.21-bookworm

WORKDIR /app

COPY . ./
RUN go mod download
RUN go build -o bin/stop-bot ./...

CMD bin/stop-bot