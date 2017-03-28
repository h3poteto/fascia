FROM golang:1.8-alpine

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

RUN chown -R go:go ${GOPATH}
USER go

# We need goose when deploy new application
# So, gom install is required
RUN set -x \
   && go get github.com/mattn/gom \
   && go get -u github.com/jteeuwen/go-bindata/... \
   && gom install \
   && go generate \
   && gom build -o bin/fascia

EXPOSE 9090:9090

CMD bin/fascia server
