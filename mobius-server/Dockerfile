# Multi-stage build for Mobius Server
FROM golang:1.23-alpine AS builder

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go.work files for multi-module workspace
COPY go.work go.work.sum ./

# Copy shared module
COPY shared/ ./shared/

# Copy all modules (required for go.work)
COPY mobius-cli/ ./mobius-cli/
COPY mobius-client/ ./mobius-client/
COPY mobius-server/ ./mobius-server/
COPY mobius-cocoon/ ./mobius-cocoon/

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
