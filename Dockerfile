# Build stage
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Create bin directory for build output
RUN mkdir -p bin

# Build the application (CGO disabled for pure Go SQLite)
# Accept VERSION build arg for versioning
ARG VERSION=dev
RUN CGO_ENABLED=0 go build -ldflags "-X main.version=${VERSION}" -o bin/mcp-email-server ./cmd/server

# Runtime stage
FROM alpine:latest

# Install CA certificates (SQLite is pure Go, no runtime deps needed)
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 mcpuser && \
    adduser -D -u 1000 -G mcpuser mcpuser

# Create data directory for SQLite cache
RUN mkdir -p /data && chown -R mcpuser:mcpuser /data

WORKDIR /app

# Copy binary from builder (updated path)
COPY --from=builder /app/bin/mcp-email-server .

# Change ownership
RUN chown mcpuser:mcpuser mcp-email-server

# Switch to non-root user
USER mcpuser

# Set environment variables
ENV CACHE_PATH=/data/email_cache.db
ENV LOG_LEVEL=info

# Run the server
CMD ["./mcp-email-server"]

