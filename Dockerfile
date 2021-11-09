FROM golang:1.16-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
COPY *.go ./

RUN go build -o /server

EXPOSE 8089
CMD ["/server"]
