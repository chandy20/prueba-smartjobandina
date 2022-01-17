#!/bin/bash

CURRENT=$(basename $PWD)

MYSQL_CONTAINER_NAME=mysql-test-$CURRENT # el nombre del contenedor con el mysql para las pruebas

DYNAMODB_CONTAINER_NAME=dynamodb-test-$CURRENT

function stopMysql(){
    # se asegura que no hayan contenedores corriendo con el nombre
    # en la variable MYSQL_CONTAINER_NAME 
    running=$(docker ps -aq -f name=$MYSQL_CONTAINER_NAME)
    if [[ ! -z "$running" ]]; then
        docker rm -f $running
    fi
}

function stopDynamoDB(){
    # se asegura que no hayan contenedores corriendo con el nombre
    # en la variable DYNAMODB_CONTAINER_NAME 
    running=$(docker ps -aq -f name=$DYNAMODB_CONTAINER_NAME)
    if [[ ! -z "$running" ]]; then
        docker rm -f $running
    fi
}

echo "**** start mysql docker"
stopMysql
docker run --rm --name $MYSQL_CONTAINER_NAME -e MYSQL_ROOT_PASSWORD=toor -e MYSQL_DATABASE=test -d mysql:5.7

echo "**** start dynamodb docker"
stopDynamoDB
docker run --rm --name $DYNAMODB_CONTAINER_NAME -d amazon/dynamodb-local:latest

# cache para el go mod
mkdir -p ~/.local-pipeline/go

echo "**** update builder image"
docker pull buildersmerqueo/builder-serverless:go1.17

echo "**** running CI"
docker run --rm -it \
    --name tester-$CURRENT \
    --link $DYNAMODB_CONTAINER_NAME \
    --link $MYSQL_CONTAINER_NAME \
    -v ~/.local-pipeline/go:/home/ubuntu/go \
    -v $PWD:/home/ubuntu/project \
    -e PACKAGES_OMIT="node_modules|01_importer|support" \
    -e CI_COVERAGE_MIN=79 \
    -e CI_STATIC_CHECK_ERRORS_MAX=0 \
    -e CI_FORMAT_ERRORS_MAX=0 \
    -e CI_LINTER_ERRORS_MAX=0 \
    -e CI_CHECK_LINTER_ERRORS=1 \
    -e DYNAMODB_URL=http://$DYNAMODB_CONTAINER_NAME:8000 \
    -e MYSQL_URL=tcp://$MYSQL_CONTAINER_NAME:3306 \
    -w /home/ubuntu/project \
    buildersmerqueo/builder-serverless:go1.17 tester-go.sh

stopMysql
stopDynamoDB