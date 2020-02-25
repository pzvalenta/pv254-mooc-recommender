# Recommender system

In this project, we implement custom recommender system.

### Installation guide

#### To run db and backend inside docker
```
docker-compose build
docker-compose up api
```

mongodb is now running on port 27017 and is restored from the dump

api is now reachable at localhost:8080/api/

try localhost:8080/api/randomRecommending

```
docker stop mongo_dev
docker stop go_api
```


#### To run only db inside docker
To run out of docker, please ensure that docker with mongo is running at localhost:27017 using the following command in parent directory:
`docker-compose up mongo_restore`.
Afterwards, you can build/run `main.go` using standard `go` commands. This should work both inside and outside of GOPATH thanks to `go.mod`.

For example, run `go run main.go`. Or if you are using Goland, just play the main function.


#### To format code
Install https://github.com/golangci/golangci-lint#install and run `golangci-lint run --fix` 
to resolve all the linting problems before every commit.

### Front end
Front end part of the project can be found [here](https://github.com/ZemanOndrej/mooc-recommender).
### Classcentral.com scraper
Scraper for our data can be found [here](https://github.com/ZemanOndrej/pv254-course-recommender-scraper)
