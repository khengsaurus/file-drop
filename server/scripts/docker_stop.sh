#!/bin/bash

docker stop $(docker ps -a | grep fd-service | cut -d " " -f 1)
docker rm $(docker ps -a | grep fd-service | cut -d " " -f 1)

docker image rm $(docker image ls | grep fd-service | tr -s ' ' | cut -d " " -f 3)