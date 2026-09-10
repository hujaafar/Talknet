# syntax=docker/dockerfile:1
FROM golang:1.26-bookworm AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY . .
FROM build AS test
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build go test -race ./... && go vet ./...
FROM build AS binary
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build CGO_ENABLED=1 go build -trimpath -ldflags="-s -w" -o /talknet .
FROM debian:bookworm-slim AS runtime
LABEL org.opencontainers.image.title="Talknet" \
      org.opencontainers.image.description="A Go and SQLite discussion community" \
      org.opencontainers.image.source="https://github.com/hujaafar/Talknet"
WORKDIR /app
RUN mkdir -p /app/data && chown -R 10001:10001 /app
COPY --from=binary --chown=10001:10001 /talknet /app/talknet
USER 10001:10001
ENV TALKNET_ADDR=:8080 TALKNET_DB=/app/data/talknet.db
EXPOSE 8080
HEALTHCHECK --interval=20s --timeout=5s --start-period=10s --retries=3 CMD ["/app/talknet", "-healthcheck"]
ENTRYPOINT ["/app/talknet"]
