# Multi-stage build for Mobius Server with Svelte UI
# Node.js stage for frontend build
FROM node:22-alpine AS frontend-builder

WORKDIR /app/frontend

# Copy frontend package files
COPY ../mobius-web/package*.json ./
COPY ../mobius-web/vite.config.ts ./
COPY ../mobius-web/svelte.config.js ./
COPY ../mobius-web/tsconfig.json ./

# Install frontend dependencies
RUN npm ci --only=production

# Copy frontend source code
COPY ../mobius-web/src ./src
COPY ../mobius-web/static ./static

# Build the frontend
RUN npm run build

# Go builder stage
FROM golang:1.24.4-alpine AS builder

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Leverage Go build cache: copy workspace files first
COPY go.work go.work.sum ./
COPY mobius-server/go.mod ./mobius-server/
COPY mobius-server/go.sum ./mobius-server/
COPY mobius-cli/go.mod ./mobius-cli/
COPY mobius-cli/go.sum ./mobius-cli/
COPY mobius-client/go.mod ./mobius-client/
COPY mobius-client/go.sum ./mobius-client/
COPY mobius-cocoon/go.mod ./mobius-cocoon/
COPY shared/go.mod ./shared/
COPY shared/go.sum ./shared/

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

# Copy frontend build output to static directory
COPY --from=frontend-builder /app/frontend/build ./mobius-server/static/

# Build the server
WORKDIR /app/mobius-server
RUN --mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o mobius-api ./cmd/api-server/main.go

# Final stage
FROM alpine:3.20

# Add non-root user and CA certs
RUN apk --no-cache add ca-certificates && \
	addgroup -S app && adduser -S -G app app

WORKDIR /app

# Copy the binary
COPY --from=builder /app/mobius-server/mobius-api /app/mobius-api
USER app

# Expose the server port
EXPOSE 8081

# Run the server
ENTRYPOINT ["/app/mobius-api"]
