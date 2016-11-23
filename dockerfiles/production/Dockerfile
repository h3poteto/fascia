FROM golang:1.7.3-alpine

ENV GOPATH /go
ENV GOJIROOT ${GOPATH}/src/github.com/h3poteto/fascia
ENV GOJIENV production

RUN set -x \
    && apk add --no-cache git mysql-client build-base \
    && mkdir -p ${GOJIROOT}

COPY . ${GOJIROOT}

WORKDIR ${GOJIROOT}

RUN set -x \
   && go get -u github.com/jteeuwen/go-bindata/... \
   && go get github.com/mattn/gom \
   && gom install \
   && go generate \
   && gom build -o bin/fascia

EXPOSE 9090:9090

CMD bin/fascia server
