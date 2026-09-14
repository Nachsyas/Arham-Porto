# Multi-stage production Dockerfile for Arham Porto Go API
# Builds server, indexer, and migrate binaries

FROM golang:1.22-alpine AS builder

WORKDIR /build

RUN apk add --no-cache ca-certificates

# Cache dependencies
COPY apps/api/go.mod apps/api/go.sum ./
RUN go mod download

# Copy source code
COPY apps/api/ ./

# Build production binaries without debug symbols
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /build/bin/server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /build/bin/indexer ./cmd/indexer
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /build/bin/migrate ./cmd/migrate

# Minimal production runtime
FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

# Copy built production binaries
COPY --from=builder /build/bin/server /app/server
COPY --from=builder /build/bin/indexer /app/indexer
COPY --from=builder /build/bin/migrate /app/migrate

# Copy canonical portfolio data and migrations
COPY data/ /app/data/
COPY apps/api/migrations/ /app/migrations/

# Production environment defaults
ENV PORT=8080
ENV PORTFOLIO_DATA_DIR=/app/data
ENV MIGRATIONS_DIR=/app/migrations
ENV APP_ENV=production
ENV TRUST_PROXY_MODE=cloudrun

EXPOSE 8080

CMD ["/app/server"]
