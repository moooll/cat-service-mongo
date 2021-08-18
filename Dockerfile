FROM golang:1.16-alpine

WORKDIR /cat-service-mongo

COPY . .
COPY ./wait-for ./wait-for
RUN ["chmod", "+x", "./wait-for"]

RUN go mod tidy
RUN go build -o cat

CMD ["sh", "-c". "./wait-for redis:6379 -- ./cat"] 