FROM h3poteto/golang-node:1.9.1 AS assets

ENV APPROOT /go/src/github.com/h3poteto/fascia
WORKDIR ${APPROOT}

USER root

COPY . ${APPROOT}

RUN chown -R go:go ${APPROOT}
USER go

RUN set -x \
    && npm install \
    && npm run release-compile


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
COPY --from=assets ${APPROOT} ${APPROOT}/public

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
