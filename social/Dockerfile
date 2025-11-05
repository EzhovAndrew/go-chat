FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /build

# Copy lib module first (shared dependency)
COPY lib/ ./lib/

# Copy service files
COPY social/go.mod social/go.sum ./social/
WORKDIR /build/social
RUN go mod download

# Copy rest of the service code
COPY social/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o social-service ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates netcat-openbsd && \
    addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

COPY --from=builder --chown=appuser:appuser /build/social/social-service .

USER appuser

EXPOSE 8080

CMD ["./social-service"]

