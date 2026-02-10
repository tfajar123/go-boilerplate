# ========================
# Build stage
# ========================
FROM golang:1.25-alpine AS builder

# Install git for go mod download and other tools
RUN apk add --no-cache git ca-certificates tzdata

# Set Go environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /app

# Copy go mod files first for better layer caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate ent code
RUN go generate ./ent

# Build the application with optimizations
RUN go build -ldflags='-w -s' -o app ./apps/cmd/server

# ========================
# Runtime stage
# ========================
FROM alpine:3.21

# Install necessary packages
RUN apk add --no-cache \
    ca-certificates \
    tzdata

# Set timezone
ENV TZ=Asia/Jakarta

WORKDIR /app

# Copy the binary and migrations
COPY --from=builder /app/app .
COPY --from=builder /app/migrations ./migrations

# Create non-root user for security
RUN adduser -D -g '' appuser && \
    chown -R appuser:appuser /app
USER appuser

EXPOSE 3098

# Use exec form for better signal handling
ENTRYPOINT ["/app/app"]
