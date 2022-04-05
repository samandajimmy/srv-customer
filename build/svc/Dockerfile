# -----------
# Build Stage
# -----------

FROM artifactory.pegadaian.co.id:8084/golang:1.17 as builder

# Set-up go proxy
ENV GO111MODULE=on \
    GOPROXY="https://artifactory.pegadaian.co.id/repository/go-group-01/"
COPY ssl_certificate.crt /usr/local/share/ca-certificates/ssl_certificate.crt

# Install some dependencies needed to build the project
RUN apt install git make \
    && mkdir /usr/src/app \
    && chmod 644 /usr/local/share/ca-certificates/ssl_certificate.crt \
    && update-ca-certificates \
    && echo ${GO111MODULE} ${GOPROXY}

# Here we copy the rest of the source code
WORKDIR /usr/src/app
COPY Makefile go.mod ./
COPY .git ./.git
COPY cmd ./cmd
COPY internal ./internal
COPY migrations ./

# Compiling...
RUN ls -la && make release

# ------------
# Deploy Stage
# ------------

FROM artifactory.pegadaian.co.id:8084/alpine:3.15

ARG ARG_PORT=3000

RUN apk add --no-cache tzdata ca-certificates

COPY --from=builder /usr/src/app/bin/release /usr/src/app

# Running...
WORKDIR /usr/src/app

EXPOSE ${ARG_PORT}

ENTRYPOINT ["./customer"]
