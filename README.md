# docker-mongo



```
docker-compose build
docker-compose up api
```

mongodb is now running on port 27017 and is restored from the dump

api is now reachable at localhost:8080/api/

```
docker stop mongo_dev
docker stop go-api
```



atm only cs courses are in the dump


to use mongo_seed:
* insert data from scraper
* `docker-compose up mongo_seed`
* to create a new dump, use mongodump
