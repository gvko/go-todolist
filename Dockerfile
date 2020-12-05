FROM golang:1.15-alpine AS build
WORKDIR /go/src/app
COPY . .
RUN go get -d -v ./...
RUN go install -v ./...

WORKDIR /app
COPY . /app
RUN go build -o main .

EXPOSE 8000
CMD ["/app/main"]
