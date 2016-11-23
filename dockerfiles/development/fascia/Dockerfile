FROM golang:1.7.3-alpine

RUN apk add --no-cache \
    git \
    build-base \
    openssh

ENV GOPATH /go
ENV GOJIROOT ${GOPATH}/src/github.com/h3poteto/fascia
RUN mkdir -p ${GOJIROOT}

WORKDIR ${GOJIROOT}

RUN go get -u github.com/mattn/gom
RUN go get -u github.com/jteeuwen/go-bindata/...

ENV GOJIENV development

EXPOSE 9090:9090

CMD ["/bin/sh"]

