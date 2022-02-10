ARG ARG_GO_VERSION=1.17

FROM golang:${ARG_GO_VERSION}-alpine

ARG ARG_SERVICE_SLUG
ARG ARG_PORT
ARG ARG_GO_MOD=go.mod

# Install some dependencies needed to build the project
RUN apk add --no-cache git build-base \
    && mkdir -p /usr/src/app

# Copy files
WORKDIR /usr/src/app
COPY configs/.air.toml /.air.toml
COPY Makefile responses.yml ./

# Building
RUN go install github.com/cosmtrek/air@v1.27.4 \
    && go install github.com/go-delve/delve/cmd/dlv@latest

EXPOSE ${ARG_PORT}

ENV SERVICE_SLUG=${ARG_SERVICE_SLUG}

CMD ["air", "-c", "/.air.toml"]
