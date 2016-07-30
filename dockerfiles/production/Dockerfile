FROM node:5.9.1-slim

ENV GOVERSION 1.6.3

RUN echo "deb http://ftp.jp.debian.org/debian/ jessie main contrib non-free" > /etc/apt/sources.list \
    && apt-get update -y \
    && apt-get install -y git mysql-client build-essential \
    && rm -rf /var/lib/apt/lists/* \
    && rm -rf /var/cache/apt/archives/* \
    && curl -O https://storage.googleapis.com/golang/go${GOVERSION}.linux-amd64.tar.gz \
    && tar -C /usr/local -xzf go${GOVERSION}.linux-amd64.tar.gz \
    && rm go${GOVERSION}.linux-amd64.tar.gz \
    && mkdir -p /go \
    && mkdir -p /root/github.com \
    && cd /root/github.com \
    && git clone https://github.com/h3poteto/fascia.git

ENV GOPATH /go
ENV GOROOT /usr/local/go
ENV PATH ${GOPATH}/bin:${GOROOT}/bin:$PATH
ENV GOJIROOT /root/github.com/fascia
ENV GO15VENDOREXPERIMENT 0
ENV GOJIENV production
WORKDIR ${GOJIROOT}

RUN go get github.com/mattn/gom \
    && gom install \
    && npm install \
    && npm run-script js-release \
    && npm run-script sass-release

EXPOSE 9090:9090

CMD gom run server.go
