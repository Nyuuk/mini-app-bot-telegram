# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Create appuser
RUN adduser -D -g '' appuser

WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

# Copy the source code
COPY . .

# Build the binary for the target architecture
RUN CGO_ENABLED=0 go build -ldflags='-w -s' -o main .

# Final stage - using alpine for better compatibility
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata wget curl

# Import user from builder
COPY --from=builder /etc/passwd /etc/passwd

# Create app directory
WORKDIR /app

# Copy the binary from builder stage
COPY --from=builder /build/main ./main

# Use an unprivileged user
USER appuser

# Expose port
EXPOSE 3000

# Add health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --quiet --tries=1 --spider http://localhost:3000/health || exit 1

# Run the binary
CMD ["./main"]
