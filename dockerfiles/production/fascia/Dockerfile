FROM golang:1.6.2-alpine

RUN apk add --no-cache git mysql-client build-base

RUN mkdir -p /go && mkdir -p /root/github.com

WORKDIR /root/github.com

RUN git clone https://github.com/h3poteto/fascia.git

ENV GOPATH /go
ENV GOJIROOT /root/github.com/fascia

ENV HOME /root
WORKDIR ${HOME}

RUN go get github.com/mattn/gom

WORKDIR ${GOJIROOT}

ENV GO15VENDOREXPERIMENT 0
ENV GOJIENV production

RUN gom install

EXPOSE 9090:9090

CMD gom run server.go
