# docker-mongo



```
docker-compose build
docker-compose up api
```

mongodb is now running on port 27017 and is restored from the dump

api is now reachable at localhost:8080/api/

try localhost:8080/api/random/ ( sadly not random right now )

```
docker stop mongo_dev
docker stop go-api
```
