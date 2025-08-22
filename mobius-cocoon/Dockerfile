# Multi-stage build for Mobius Cocoon
FROM golang:1.24.4-alpine AS builder

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Prepare cached deps
COPY go.work go.work.sum ./
COPY mobius-cocoon/go.mod mobius-cocoon/go.sum ./mobius-cocoon/
COPY mobius-server/go.mod mobius-server/go.sum ./mobius-server/
COPY mobius-cli/go.mod mobius-cli/go.sum ./mobius-cli/
COPY mobius-client/go.mod mobius-client/go.sum ./mobius-client/
COPY shared/go.mod shared/go.sum ./shared/
RUN --mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	go work sync && \
	(cd mobius-cocoon && go mod download) && \
	(cd mobius-server && go mod download) && \
	(cd mobius-cli && go mod download) && \
	(cd mobius-client && go mod download) && \
	(cd shared && go mod download)

# Copy sources
COPY . .

# Build the cocoon
WORKDIR /app/mobius-cocoon
RUN --mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o cocoon ./cmd/cocoon

# Final stage
FROM alpine:3.20
RUN apk --no-cache add ca-certificates && \
	addgroup -S app && adduser -S -G app app
WORKDIR /app
COPY --from=builder /app/mobius-cocoon/cocoon /app/cocoon
USER app

# Expose the cocoon port
EXPOSE 8082

# Run the cocoon
ENTRYPOINT ["/app/cocoon"]
