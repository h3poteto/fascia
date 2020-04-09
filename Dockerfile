FROM node:10.16.3-alpine AS assets

ENV APPROOT /var/opt/app

COPY . ${APPROOT}

WORKDIR ${APPROOT}/assets

RUN set -x \
    && npm install \
    && npm run compile

FROM node:10.16.3-alpine AS lp

ENV APPROOT /var/opt/app

COPY . ${APPROOT}

WORKDIR ${APPROOT}/lp

RUN set -x \
    && npm install \
    && npm run compile


FROM h3poteto/golang:1.13.4

USER root
ENV GOPATH /go
ENV APPROOT ${GOPATH}/src/github.com/h3poteto/fascia
ENV APPENV production
ENV GO111MODULE on

RUN set -x \
    && apk add --no-cache \
    curl && \
    curl -fsSL https://github.com/minamijoyo/myaws/releases/download/v0.3.0/myaws_v0.3.0_linux_amd64.tar.gz \
    | tar -xzC /usr/local/bin && chmod +x /usr/local/bin/myaws

WORKDIR ${APPROOT}

COPY --chown=go:go . ${APPROOT}
COPY --chown=go:go --from=assets /var/opt/app/public/assets ${APPROOT}/public/assets
COPY --chown=go:go --from=lp /var/opt/app/public/lp ${APPROOT}/public/lp

RUN chown -R go:go ${APPROOT}

USER go

RUN set -x \
   && go mod download \
   && go generate \
   && go build -o bin/fascia

EXPOSE 9090:9090

ENTRYPOINT ["./entrypoint.sh"]

CMD bin/fascia server
