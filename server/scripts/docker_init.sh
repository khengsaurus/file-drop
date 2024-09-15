#!/bin/bash

COMMIT_HASH=$(git rev-parse --short HEAD)

FULL_APP_NAME="fd-service"
FULL_IMAGE_NAME="$FULL_APP_NAME:$COMMIT_HASH"
if [[ "$(docker images -q $FULL_IMAGE_NAME 2> /dev/null)" == "" ]]; 
  then docker build . -t $FULL_IMAGE_NAME --build-arg service=all
fi

docker run --name $FULL_APP_NAME-API -p "8094:8080" -d $FULL_IMAGE_NAME
docker run --name $FULL_APP_NAME-STR -p "8095:8080" -d $FULL_IMAGE_NAME
docker run --name $FULL_APP_NAME-URL-1 -p "8096:8080" -d $FULL_IMAGE_NAME
docker run --name $FULL_APP_NAME-URL-1 -p "8097:8080" -d $FULL_IMAGE_NAME

# Build app image for API service if not exists
# API_APP_NAME="fd-service-api"
# API_IMAGE_NAME="$API_APP_NAME:$COMMIT_HASH"
# if [[ "$(docker images -q $API_IMAGE_NAME 2> /dev/null)" == "" ]]; 
#   then docker build . -t $API_IMAGE_NAME --build-arg service=api
# fi

# docker run --name $API_APP_NAME-1 -p "8094:8080" -d $API_IMAGE_NAME

# # Build app image for Stream service if not exists
# STR_APP_NAME="fd-service-str"
# STR_IMAGE_NAME="$STR_APP_NAME:$COMMIT_HASH"
# if [[ "$(docker images -q $STR_IMAGE_NAME 2> /dev/null)" == "" ]]; 
#   then docker build . -t $STR_IMAGE_NAME --build-arg service=stream
# fi

# docker run --name $STR_APP_NAME-1 -p "8095:8080" -d $STR_IMAGE_NAME

# # Build app image for URL service if not exists
# URL_APP_NAME="fd-service-url"
# URL_IMAGE_NAME="$URL_APP_NAME:$COMMIT_HASH"
# if [[ "$(docker images -q $URL_IMAGE_NAME 2> /dev/null)" == "" ]];
#   then docker build . -t $URL_IMAGE_NAME --build-arg service=url
# fi

# docker run --name $URL_APP_NAME-1 -p "8096:8080" -d $URL_IMAGE_NAME
# docker run --name $URL_APP_NAME-2 -p "8097:8080" -d $URL_IMAGE_NAME