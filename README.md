# simple-server
This is a simple server implemented in Go.

## Setup
Start by building the docker images with
```sh
docker build -t postgres-example -f Dockerfile.postgres .
docker build -t simple-server .
```

Then start the database with
```sh
docker run --name postgres-db -d -p 5432:5432 postgres-example
```
and start the server with
```sh
docker run --name simple-server -d -p 8080:8080 server
```
