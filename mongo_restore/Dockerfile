FROM mongo:latest

COPY dump /dump

CMD mongorestore --drop --host=mongo dump/

