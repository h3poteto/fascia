FROM golang:1.6.2-alpine

RUN apk add --no-cache git mysql-client build-base

RUN mkdir -p /go && mkdir -p /root/fascia

ENV GOPATH /go
ENV GOJIROOT /root/fascia

ENV HOME /root
WORKDIR ${HOME}

RUN go get github.com/mattn/gom

COPY . /root/fascia

WORKDIR /root/fascia

ENV GO15VENDOREXPERIMENT 0
ENV GOJIENV production

RUN gom install

EXPOSE 9090:9090

CMD gom run server.go
