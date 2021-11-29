FROM node:14.18-slim AS assets

ENV APPROOT /var/opt/app

COPY . ${APPROOT}

WORKDIR ${APPROOT}/assets

RUN set -x \
    && npm install \
    && npm run compile

FROM node:14.18-slim AS lp

ENV APPROOT /var/opt/app

COPY . ${APPROOT}

WORKDIR ${APPROOT}/lp

RUN set -x \
    && npm install \
    && npm run compile


FROM ghcr.io/h3poteto/golang:1.16.10

USER root
ENV GOPATH /go
ENV APPROOT ${GOPATH}/src/github.com/h3poteto/fascia
ENV APPENV production
ENV GO111MODULE on

WORKDIR ${APPROOT}

COPY --chown=go:go . ${APPROOT}
COPY --chown=go:go --from=assets /var/opt/app/public/assets ${APPROOT}/public/assets
COPY --chown=go:go --from=lp /var/opt/app/public/lp ${APPROOT}/public/lp

RUN chown -R go:go ${APPROOT}

USER go

RUN set -x \
   && go mod download \
   && go build -o bin/fascia

EXPOSE 9090:9090

ENTRYPOINT ["./entrypoint.sh"]

CMD bin/fascia server
