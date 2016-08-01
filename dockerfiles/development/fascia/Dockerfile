FROM golang:1.6.2-alpine

RUN apk add --no-cache \
    git \
    build-base \
    openssh

RUN mkdir -p /go && \
    mkdir -p /root/fascia

ENV GOPATH /go
ENV GOJIROOT /root/fascia

ENV HOME /root
WORKDIR ${HOME}

RUN go get github.com/mattn/gom

ENV GO15VENDOREXPERIMENT 0
ENV GOJIENV development

EXPOSE 9090:9090

WORKDIR /root/fascia

CMD ["/bin/sh"]

