#!/bin/bash

# Init app as a docker image, using remote services

APP="fd-server"
TAG="1"
PORTS="8094:8080"
BUILD_PATH="."

# Build app image if not exists
if [[ "$(docker images -q $APP:$TAG 2> /dev/null)" == "" ]]; 
  then docker build $BUILD_PATH -t $APP:$TAG
fi

docker run --name $APP -p $PORTS -e APP_ID=$APP -d $APP:$TAG