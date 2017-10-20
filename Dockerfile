FROM h3poteto/golang:1.9.1

USER root
ENV GOPATH /go
ENV APPROOT ${GOPATH}/src/github.com/h3poteto/fascia
ENV APPENV production

RUN set -x \
    && apk add --no-cache \
    curl

WORKDIR ${APPROOT}

COPY . ${APPROOT}

RUN chown -R go:go ${GOPATH}
USER go

RUN set -x \
   && go get github.com/mattn/gom \
   && go get -u github.com/jteeuwen/go-bindata/... \
   && dep ensure \
   && go generate \
   && go build -o bin/fascia

EXPOSE 9090:9090

CMD bin/fascia server
