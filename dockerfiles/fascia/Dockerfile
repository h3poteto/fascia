FROM golang:1.6.2-alpine

RUN apk add --no-cache git

RUN mkdir -p /go
RUN mkdir -p /root/fascia

ENV GOPATH /go
ENV GOJIROOT /root/fascia

ENV HOME /root
WORKDIR ${HOME}

RUN go get github.com/mattn/gom

EXPOSE 9090:9090

WORKDIR /root/fascia

CMD ["/bin/sh"]

