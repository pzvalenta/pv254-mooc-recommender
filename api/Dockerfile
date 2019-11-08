FROM golang:alpine as build-env

ENV GO111MODULE=on \
    GOPATH=/go \
    # DB_HOST=localhost \
    # DB_PORT=5432 \
    # DB_USER=docker \
    # DB_PW=docker \
    # DB_NAME=docker \
    CGO_ENABLED=0 

RUN mkdir -p /go/src/go-docker-api/api

WORKDIR /go/src/go-docker-api/api

COPY ./go.mod . 
COPY ./go.sum .
RUN apk add --no-cache git  \
    && go mod download \
    && apk del git 

COPY . .

RUN go build -a -installsuffix cgo -o /go/src/go-docker-api/api/main

# FROM scratch 
# COPY --from=build-env /go/src/go-docker-api/api/main /go/src/go-docker-api/api/src/main

ENTRYPOINT ["/go/src/go-docker-api/api/main"]
