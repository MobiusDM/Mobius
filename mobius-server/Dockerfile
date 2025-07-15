# Multi-stage build for Mobius Server
FROM golang:1.24.4-alpine AS builder

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go.work files for multi-module workspace
COPY go.work go.work.sum ./

# Copy shared module
COPY shared/ ./shared/

# Copy server module
COPY mobius-server/ ./mobius-server/

# Build the server
WORKDIR /app/mobius-server
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mobius ./cmd/mobius

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary
COPY --from=builder /app/mobius-server/mobius .

# Expose the server port
EXPOSE 8080

# Run the server
CMD ["./mobius"]
