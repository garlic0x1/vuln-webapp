#!/bin/bash

docker build -t "srv" src/
docker build -t "mysql" db/

# try this
set +e
# create my-net if not already exists
docker network create my-net

# remove these if existing
docker container stop /srv
docker container rm /srv
docker container stop /mysql
docker container rm /mysql
set -e

docker run --name mysql --network my-net -d -p 3306:3306 mysql
sleep 10
docker run --name srv --network my-net -d -p 9001:80 srv

