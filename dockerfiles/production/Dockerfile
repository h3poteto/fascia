FROM golang:1.6.2-alpine

ENV GOPATH /go
ENV GOJIROOT /root/fascia
ENV GO15VENDOREXPERIMENT 0
ENV GOJIENV production

RUN set -x \
    && apk add --no-cache git mysql-client build-base \
    && mkdir -p /go \
    && mkdir -p ${GOJIROOT}

COPY . ${GOJIROOT}

WORKDIR ${GOJIROOT}

RUN set -x \
   && go get -u github.com/jteeuwen/go-bindata/... \
   && go get github.com/mattn/gom \
   && gom install \
   && go generate \
   && gom build -o bin/fascia server.go

EXPOSE 9090:9090

CMD bin/fascia
