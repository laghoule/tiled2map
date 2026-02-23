##############################
FROM golang:1.26-alpine AS build

ARG VERSION "devel"
ARG GIT_COMMIT ""

WORKDIR /src

RUN --mount=type=bind,source=.,target=.  \
  --mount=type=cache,target=/root/.cache/go-build \
  --mount=type=cache,target=/go/pkg \
  CGO_ENABLED=0 GOOS=linux go build cmd/tile2map.go -ldflags="-X 'main.version=$VERSION' -X 'main.gitCommit=$GIT_COMMIT'" -o /tmp/tiled2map main.go

##############################
FROM scratch

ARG VERSION

LABEL org.opencontainers.image.title="tiled2map" \
  org.opencontainers.image.vendor="laghoule" \
  org.opencontainers.image.licenses="GPLv3" \
  org.opencontainers.image.version="${VERSION}" \
  org.opencontainers.image.description="TODO!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!." \
  org.opencontainers.image.url="https://github.com/laghoule/tiled2map/README.md" \
  org.opencontainers.image.source="https://github.com/laghoule/tiled2map" \
  org.opencontainers.image.documentation="https://github.com/laghoule/tiled2map/README.md"

USER 65534

COPY --link --from=build /tmp/tiled2map /bin/tiled2map

ENTRYPOINT ["/bin/tiled2map"]
