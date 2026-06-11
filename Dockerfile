FROM golang:1.26-alpine AS builder

# Multi-repo build: the build context must be the parent directory holding both
# warden-engine and warden-auth, because go.mod has
# `replace github.com/NicolasPaterno/warden-auth => ../warden-auth` (sibling repo).
# Build with context at the parent, e.g.:
#   docker build -f warden-engine/Dockerfile -t warden-engine ..
WORKDIR /src

COPY warden-auth/go.mod warden-auth/go.sum ./warden-auth/
COPY warden-engine/go.mod warden-engine/go.sum ./warden-engine/
WORKDIR /src/warden-engine
RUN go mod download

WORKDIR /src
COPY warden-auth/ ./warden-auth/
COPY warden-engine/ ./warden-engine/

WORKDIR /src/warden-engine
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /engine ./cmd/engine

# ─────────────────────────────────────────────────────────────────────────────

FROM alpine:3.21

RUN addgroup -S warden && adduser -S engine -G warden

WORKDIR /app

COPY --from=builder /engine .

USER engine

EXPOSE 8081

ENTRYPOINT ["/app/engine"]
