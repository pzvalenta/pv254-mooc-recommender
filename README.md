docker-mongo



docker-compose build
docker-compose up


mongodb is now running on port 27017 and is restored from the dump
atm only cs courses are in the dump

to create a new dump, use mongodump


to use mongo_seed:
-  insert data from scraper
-  uncomment the section in docker-compose.yml
