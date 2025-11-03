# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /build/gateway

# Layer 1: Copy dependency files (rarely changes - good for caching)
# Gateway needs access to other services' go.mod for replace directives
COPY auth/go.mod auth/go.sum ../auth/
COPY users/go.mod users/go.sum ../users/
COPY chat/go.mod chat/go.sum ../chat/
COPY social/go.mod social/go.sum ../social/
COPY notifications/go.mod notifications/go.sum ../notifications/
COPY gateway/go.mod gateway/go.sum ./

# Layer 2: Download dependencies (cached if go.mod/go.sum unchanged)
RUN go mod download

# Layer 3: Copy all source code (changes frequently - at the end)
COPY auth/pkg ../auth/pkg
COPY users/pkg ../users/pkg
COPY chat/pkg ../chat/pkg
COPY social/pkg ../social/pkg
COPY notifications/pkg ../notifications/pkg
COPY gateway .

# Layer 4: Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gateway ./cmd/main.go

# ============================================================================
# Runtime stage - minimal image
# ============================================================================
FROM alpine:latest

# Single layer for all setup (combines apk + user creation)
RUN apk --no-cache add ca-certificates wget && \
    addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Copy binary and set ownership in one layer
COPY --from=builder --chown=appuser:appuser /build/gateway/gateway .

USER appuser

EXPOSE 8080

CMD ["./gateway"]

