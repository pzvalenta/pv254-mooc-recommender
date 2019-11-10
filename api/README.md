# api

To run out of docker, please ensure that docker with mongo is running at localhost:27017 using the following command in parent directory: `docker-compose up mongo_seed`

Afterwards, you can build/run `main.go` using standard `go` commands. This should work both inside and outside of GOPATH thanks to `go.mod`.

