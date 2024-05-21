FROM --platform=$BUILDPLATFORM golang:alpine AS build-server
WORKDIR /AV
COPY . /AV
ARG TARGETOS
ARG TARGETARCH
ENV CGO_ENABLED=0
RUN set -ex \
    && GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o AV -tags "sing xray with_reality_server with_quic with_grpc with_utls with_wireguard with_acme" -trimpath -ldflags "-s -w"

FROM --platform=${TARGETPLATFORM} alpine:latest
WORKDIR /AV
COPY --from=build-server /AV/AV /AV

ENV LANG C.UTF-8
ENTRYPOINT [ "/AV/AV", "server", "--config", "/etc/V2bX/config.json" ]
