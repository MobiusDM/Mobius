# Go builder stage
FROM golang:1.24-alpine AS builder

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go workspace files
COPY go.work go.work.sum ./

# Copy module files
COPY mobius-server/go.mod mobius-server/go.sum ./mobius-server/
COPY mobius-shared/go.mod ./mobius-shared/

# Download dependencies
WORKDIR /app/mobius-server
RUN go mod download

# Back to app directory
WORKDIR /app

# Copy source code
COPY mobius-server/ ./mobius-server/
COPY mobius-shared/ ./mobius-shared/

# Create minimal assets directory for bindata
RUN mkdir -p mobius-server/assets && echo "/* Empty CSS */" > mobius-server/assets/bundle.css

# Generate embedded assets and build the application
WORKDIR /app/mobius-server
RUN go run github.com/kevinburke/go-bindata/go-bindata -pkg=bindata -tags full \
    -o=server/bindata/generated.go \
    assets/... server/mail/templates

RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -installsuffix cgo -o mobius ./cmd/mobius

# Production stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/mobius .

# Expose port 8080
EXPOSE 8080

# Run the application
CMD ["./mobius"]