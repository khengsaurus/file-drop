#!/bin/bash

# Build app image for API service if not exists
APP_API="fd-service-api"
TAG_API="1_api"
if [[ "$(docker images -q $APP_API:$TAG_API 2> /dev/null)" == "" ]]; 
  then docker build . -t $APP_API:$TAG_API --build-arg service=api
fi

docker run --name $APP_API-1 -p "8094:8080" -d $APP_API:$TAG_API

# Build app image for Stream service if not exists
APP_STR="fd-service-stream"
TAG_STR="1_str"
if [[ "$(docker images -q $APP_STR:$TAG_STR 2> /dev/null)" == "" ]]; 
  then docker build . -t $APP_STR:$TAG_STR --build-arg service=stream
fi

docker run --name $APP_STR-1 -p "8095:8080" -d $APP_STR:$TAG_STR

# Build app image for URL service if not exists
APP_URL="fd-service-url"
TAG_URL="1_url"
if [[ "$(docker images -q $APP_URL:$TAG_URL 2> /dev/null)" == "" ]];
  then docker build . -t $APP_URL:$TAG_URL --build-arg service=url
fi

docker run --name $APP_URL-1 -p "8096:8080" -d $APP_URL:$TAG_URL
docker run --name $APP_URL-2 -p "8097:8080" -d $APP_URL:$TAG_URL