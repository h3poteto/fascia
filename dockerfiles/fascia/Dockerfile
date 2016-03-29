FROM golang:1.5

RUN apt-get update && apt-get -y install git mysql-client

RUN mkdir -p /go/src/fascia && mkdir -p /go/bin

ENV GOPATH /go
ENV GOJIROOT /go/src/fascia

ENV HOME /root
WORKDIR ${HOME}
RUN git clone http://github.com/zimbatm/direnv
RUN cd direnv && make install
RUN echo 'eval "$(direnv hook bash)"' >> ${HOME}/.bashrc

RUN go get github.com/mattn/gom

EXPOSE 9090:9090

WORKDIR /go/src/fascia

CMD ["/bin/bash"]

