# simple-server
This is a simple server implemented in Go.

## Setup
Start with
```sh
docker-compose up -d
```
This will build and pull the corresponding images.


## Usage
Request data from the server by executing one of the following calls:
```sh
curl http://localhost:8080/city/{cityname}
curl http://localhost:8080/city
curl http://localhost:8080/city?limit=10
```
