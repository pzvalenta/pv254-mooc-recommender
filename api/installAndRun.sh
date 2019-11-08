#!/bin/bash
PORT="8080"
INNER_PORT="8080"
while getopts p:i: option
do
case "${option}"
in
p) PORT="${OPTARG}";;
i) INNER_PORT="${OPTARG}";;
esac
done

docker build -t golang-api .
docker run --rm -p ${PORT}:${INNER_PORT} --env PORT=${INNER_PORT} golang-api
echo "docker container successfully created with ${INNER_PORT} mapped to ${PORT}"