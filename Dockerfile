FROM golang:1.21-alpine
RUN apk update
RUN apk upgrade
RUN apk add --no-cache gcc
RUN apk add --no-cache make 
RUN apk add --no-cache musl-dev
RUN apk add --no-cache mpg123-dev
RUN apk add --no-cache opus-dev

WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -gcflags "-N -l" -v -o /usr/local/bin/app ./...

CMD ["app"]