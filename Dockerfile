FROM golang:1.26-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o engine ./cmd/engine/main.go

# ─────────────────────────────────────────────────────────────────────────────

FROM alpine:3.21

RUN addgroup -S warden && adduser -S engine -G warden

WORKDIR /app

COPY --from=builder /app/engine .

USER engine

EXPOSE 8081

ENTRYPOINT ["/app/engine"]
