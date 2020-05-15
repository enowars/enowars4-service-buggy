FROM golang:latest

RUN apt-get update && apt-get install -y make gcc

ENV GOPATH /go/src

EXPOSE 7890
RUN go get github.com/gorilla/mux
RUN go get github.com/gorilla/sessions

WORKDIR /service/buggy

CMD make run
