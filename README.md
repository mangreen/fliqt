# Fli-Qt: Golang Take Home Assignment
https://www.fliqt.io/

## Description
Write Go GIN RESTful API Service include below Condition (Push to Github)
・HR System Backend API
・MySQL as Database
・Redis as Cache
・GORM Migration
・GORM MySQL SEED data
・Unit Test
・Makefile for build and deploy
・Run all service using docker-compose

## Run
```shell
# run depend services
docker-compose up -d

# run server
go run cmd/main.go
```
Web Server will run on http://localhost:8080

## Test
```shell
# all
go test -v ./...

# single
go test -v ./{file}
```

## Docker
```shell
# build
docker build -t fliqt .

# run
docker run --rm -p 8080:8080 --name fliqt fliqt
```