FROM golang:1.24-alpine AS build

# can be passed with any prefix (like `v1.2.3@FOO`), e.g.: `docker build --build-arg "APP_VERSION=v1.2.3@FOO" .`
ARG APP_VERSION="undefined@docker"

WORKDIR /src

# copy the source code
COPY . /src

RUN CGO_ENABLED=0 go build \
      -trimpath \
      -ldflags "-s -w -X github.com/micronull/pocketbook-cloud-sync/internal/pkg/version.version=${APP_VERSION}" \
      -o ./pbcsync \
      ./cmd/pbcsync/ \
    && ./pbcsync version  \
    && mkdir -p /tmp/rootfs \
    && cd /tmp/rootfs \
    && mkdir -p ./etc/ssl/certs ./bin ./tmp ./books \
    && echo 'appuser:x:10001:10001::/nonexistent:/sbin/nologin' > ./etc/passwd \
    && echo 'appuser:x:10001:' > ./etc/group \
    && chmod 777 ./tmp ./books \
    && cp /etc/ssl/certs/ca-certificates.crt ./etc/ssl/certs/ \
    && mv /src/pbcsync ./bin/pbcsync

FROM scratch AS runtime

ARG APP_VERSION="undefined@docker"

ENV DIR="/books"

LABEL \
    # Docs: <https://github.com/opencontainers/image-spec/blob/master/annotations.md>
    org.opencontainers.image.title="pbcsync" \
    org.opencontainers.image.description="Download you PocketBook Cloud library" \
    org.opencontainers.image.url="https://github.com/micronull/pocketbook-cloud-sync" \
    org.opencontainers.image.source="https://github.com/micronull/pocketbook-cloud-sync" \
    org.opencontainers.image.vendor="micronull" \
    org.opencontainers.version="$APP_VERSION" \
    org.opencontainers.image.licenses="MIT"

# import compiled application
COPY --from=build /tmp/rootfs /

# use an unprivileged user
USER 10001:10001

ENTRYPOINT ["/bin/pbcsync"]

CMD ["sync", "-env"]