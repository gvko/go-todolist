FROM golang:1.15-alpine

WORKDIR /app
ADD . /app

RUN apk add git \
  && go get -d -v ./... \
  && go build -o main .

CMD ["/app/main"]
