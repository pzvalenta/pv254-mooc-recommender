FROM golang:alpine as build-env

ENV GO111MODULE=on \
    GOPATH=/go \
    # DB_HOST=localhost \
    # DB_PORT=5432 \
    # DB_USER=docker \
    # DB_PW=docker \
    # DB_NAME=docker \
    CGO_ENABLED=0 

RUN mkdir -p /go/src/github.com/ZaxCZ/docker-mongo/api

WORKDIR /go/src/github.com/ZaxCZ/docker-mongo/api

COPY ./go.mod . 
COPY ./go.sum .
RUN apk add --no-cache git  \
    && go mod download \
    && apk del git

COPY . .

RUN go build -o main

# FROM scratch 
# COPY --from=build-env /go/src/github.com/ZaxCZ/docker-mongo/api/main /go/src/github.com/ZaxCZ/docker-mongo/api/src/main

ENTRYPOINT ["/go/src/github.com/ZaxCZ/docker-mongo/api/main"]
