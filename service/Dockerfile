FROM golang:latest

RUN apt-get update && apt-get install -y make gcc

ENV GOPATH /go/src

EXPOSE 7890
RUN go get github.com/gorilla/mux
RUN go get github.com/gorilla/sessions
RUN go get github.com/go-sql-driver/mysql

WORKDIR /service/buggy

CMD make run
