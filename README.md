#docker-mongo



```
docker-compose build
docker-compose up mongo_restore
```

mongodb is now running on port 27017 and is restored from the dump


```
docker stop mongo_dev
```



atm only cs courses are in the dump



to use mongo_seed:
* insert data from scraper
* `docker-compose up mongo_seed`
* to create a new dump, use mongodump
