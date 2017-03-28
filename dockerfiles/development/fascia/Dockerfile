FROM golang:1.8-alpine

RUN set -x \
    && apk add --no-cache \
    git \
    build-base \
    mysql-client

ENV GOPATH /go
ENV GOJIROOT ${GOPATH}/src/github.com/h3poteto/fascia
ENV GOJIENV development

RUN set -x \
    && adduser -D -u 1000 go \
    && echo 'go:password' | chpasswd \
    && mkdir -p ${GOJIROOT} \
    && chown -R go:go /go

WORKDIR ${GOJIROOT}
USER go

RUN go get -u github.com/mattn/gom
RUN go get -u github.com/jteeuwen/go-bindata/...

EXPOSE 9090:9090

CMD ["/bin/sh"]

