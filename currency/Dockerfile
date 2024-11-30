FROM golang:alpine3.20

## Install git
# RUN apk add git

RUN mkdir /app
ADD . /app
WORKDIR /app

RUN go mod download
RUN go build -o main .

ENV GIN_MODE "release"
CMD ["/app/main"]