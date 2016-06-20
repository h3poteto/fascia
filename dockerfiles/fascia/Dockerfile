FROM golang:1.6.2

RUN apt-get update && apt-get -y install git mysql-client

RUN mkdir -p /go
RUN mkdir -p /root/fascia

ENV GOPATH /go
ENV GOJIROOT /root/fascia

ENV HOME /root
WORKDIR ${HOME}

RUN go get github.com/mattn/gom

EXPOSE 9090:9090

WORKDIR /root/fascia

CMD ["/bin/bash"]

