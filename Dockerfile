FROM golang:1.16-buster

WORKDIR /cat-service-mongo

COPY . .
COPY ./wait-for ./wait-for
RUN ["chmod", "+x", "./wait-for"]
RUN apt-get update && apt-get install -y netcat

RUN go mod vendor
RUN go build -o cat

