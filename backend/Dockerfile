# Build stage
FROM golang:1.24-alpine AS builder

# Install git and ca-certificates
RUN apk update && apk add --no-cache git ca-certificates tzdata && update-ca-certificates

# Create appuser
RUN adduser -D -g '' appuser

WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod go.sum ./
RUN go mod download
RUN go mod verify

# Copy the source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' -a \
    -installsuffix cgo -o main .

# Final stage
FROM scratch

# Import from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/passwd /etc/passwd

# Copy our static executable
COPY --from=builder /build/main /go/bin/main

# Use an unprivileged user
USER appuser

# Expose port
EXPOSE 3000

# Run the binary
ENTRYPOINT ["/go/bin/main"]
