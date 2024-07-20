## Description

Hai! this is my first ever API using Go so have fun to check out and i know it has a lot of flaws so let me know how to improve it that will be a huge help. With this one you also can deploy it using docker so have fun try to using it!. By default it will be running on localhost:6012 unless you've change it, further more kindly check `/api-docs` for api documentation.

PS: in this API documentation i don't know how to change OpenAPI configuration to use Bearer Auth so i suggest u to test it useing post man with deatil taht you got form `/api-docs`

## Installation

1. Please check `.env.example` file for database connection and JWT secret then delete `.example` from the filename.
2. I recommend you to install air to user watchmode like nodemon by `go install github.com/air-verse/air@latest`

```bash
# installing dependencies
$ go mod download

# build Go app
$ go build -o main.go
```

## Running the app

```bash
# watch mode
$ air

# production mode
$ go run main.go
or
Run the binary file
```

## Deployment using docker

```bash
# build an image
$ docker-compose build

# running container
$ docker-compose up

# running container on background
$ docker-compose up -d
```
