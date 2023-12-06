# URL Shortener

## Overview

This app allows to shorten long links.

## Prerequisites:

Using redis as data storage, so need to install redis or run docker image for redis
```sh
docker run -d -p 6379:6379 --name redis redis
```

## Environment variables
SHORT_BASE_URL= <Host Name><br/>
PORT= <Port><br/>
DB_ADDR= <Redis address - localhost:6379><br/>
DB_PASS= < Redis Password><br/>

## Steps

1. Create an .env file for set env variables
2. Enter the details of environment variables in .env file
3.  Run the code with the command: go run cmd/main.go

## Steps for build docker image with Dockerfile

1. Go to root folder 
2.  Run the commands:
```sh
docker build -t <image_name> -f ./docker/Dockerfile .
```

