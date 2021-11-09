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

## Usage
Request data from the server by executing one of the following calls:
```sh
curl http://localhost:8080/city/{cityname}
curl http://localhost:8080/city
curl http://localhost:8080/city?limit=10
```
