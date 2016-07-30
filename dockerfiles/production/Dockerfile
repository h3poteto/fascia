FROM golang:1.6.2-alpine

ENV GOPATH /go
ENV GOJIROOT /root/fascia
ENV GO15VENDOREXPERIMENT 0
ENV GOJIENV production

RUN apk add --no-cache git mysql-client build-base \
    && mkdir -p /go \
    && mkdir -p ${GOJIROOT}

COPY . ${GOJIROOT}

WORKDIR ${GOJIROOT}

RUN go get github.com/mattn/gom \
   && gom install

EXPOSE 9090:9090

CMD gom run server.go
