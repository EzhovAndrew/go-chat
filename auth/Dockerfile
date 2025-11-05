FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /build

# Copy lib module first (shared dependency)
COPY lib/ ./lib/

# Copy service files
COPY auth/go.mod auth/go.sum ./auth/
WORKDIR /build/auth
RUN go mod download

# Copy rest of the service code
COPY auth/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auth-service ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates netcat-openbsd && \
    addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

COPY --from=builder --chown=appuser:appuser /build/auth/auth-service .

USER appuser

EXPOSE 8080

CMD ["./auth-service"]

