# Multi-stage build for Mobius Server
FROM golang:1.24.4-alpine AS builder

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Leverage Go build cache: copy workspace files first
COPY go.work go.work.sum ./
COPY mobius-server/go.mod mobius-server/go.sum ./mobius-server/
COPY mobius-cli/go.mod mobius-cli/go.sum ./mobius-cli/
COPY mobius-client/go.mod mobius-client/go.sum ./mobius-client/
COPY mobius-cocoon/go.mod mobius-cocoon/go.sum ./mobius-cocoon/
COPY shared/go.mod shared/go.sum ./shared/

# Download deps (cached)
RUN --mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	go work sync && \
	(cd mobius-server && go mod download) && \
	(cd mobius-cli && go mod download) && \
	(cd mobius-client && go mod download) && \
	(cd mobius-cocoon && go mod download) && \
	(cd shared && go mod download)

# Copy the full sources
COPY . .

# Build the server
WORKDIR /app/mobius-server
RUN --mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o mobius ./cmd/mobius

# Final stage
FROM alpine:3.20

# Add non-root user and CA certs
RUN apk --no-cache add ca-certificates && \
	addgroup -S app && adduser -S -G app app

WORKDIR /app

# Copy the binary
COPY --from=builder /app/mobius-server/mobius /app/mobius
USER app

# Expose the server port
EXPOSE 8080

# Run the server
ENTRYPOINT ["/app/mobius"]
