FROM golang:1.7.4-alpine

RUN set -x \
    && apk add --no-cache \
    git \
    mysql-client \
    build-base

ENV GOPATH /go
ENV GOJIROOT ${GOPATH}/src/github.com/h3poteto/fascia
ENV GOJIENV production

RUN set -x \
    && adduser -D -u 1000 go \
    && echo 'go:password' | chpasswd \
    && mkdir -p ${GOJIROOT}

WORKDIR ${GOJIROOT}

COPY . ${GOJIROOT}

RUN chown -R go:go /go
USER go

RUN set -x \
   && go get -u github.com/jteeuwen/go-bindata/... \
   && go get github.com/mattn/gom \
   && gom install \
   && go generate \
   && gom build -o bin/fascia

EXPOSE 9090:9090

CMD bin/fascia server
