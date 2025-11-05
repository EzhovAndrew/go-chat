FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git

WORKDIR /build

# Copy lib module first (shared dependency)
COPY lib/ ./lib/

# Copy service files
COPY chat/go.mod chat/go.sum ./chat/
WORKDIR /build/chat
RUN go mod download

# Copy rest of the service code
COPY chat/ ./

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o chat-service ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates netcat-openbsd && \
    addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser

WORKDIR /app

COPY --from=builder --chown=appuser:appuser /build/chat/chat-service .

USER appuser

EXPOSE 8080

CMD ["./chat-service"]

