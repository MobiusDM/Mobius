# Go builder stage
FROM golang:1.24-alpine AS builder

# Install git and ca-certificates
RUN apk add --no-cache git ca-certificates

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy backend source code
COPY backend/ ./backend/

# Copy server mail templates
COPY backend/server/mail/templates ./server/mail/templates

# Create minimal assets directory for bindata
RUN mkdir -p assets && echo "/* Empty CSS */" > assets/bundle.css

# Generate embedded assets and build the application
RUN go run github.com/kevinburke/go-bindata/go-bindata -pkg=bindata -tags full \
    -o=backend/server/bindata/generated.go \
    assets/... server/mail/templates

RUN cd backend && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ../mobius ./cmd/mobius

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