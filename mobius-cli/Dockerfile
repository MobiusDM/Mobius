# Multi-stage build for Mobius CLI
FROM golang:1.24.4-alpine AS builder

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Leverage Go build cache
COPY go.work go.work.sum ./
COPY mobius-cli/go.mod ./mobius-cli/
COPY mobius-cli/go.sum ./mobius-cli/
COPY mobius-client/go.mod ./mobius-client/
COPY mobius-client/go.sum ./mobius-client/
COPY mobius-server/go.mod ./mobius-server/
COPY mobius-server/go.sum ./mobius-server/
COPY mobius-cocoon/go.mod ./mobius-cocoon/
COPY shared/go.mod ./shared/
COPY shared/go.sum ./shared/
RUN --mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	go work sync && \
	(cd mobius-cli && go mod download) && \
	(cd mobius-client && go mod download) && \
	(cd mobius-server && go mod download) && \
	(cd mobius-cocoon && go mod download) && \
	(cd shared && go mod download)

# Copy source
COPY . .

# Build the CLI
WORKDIR /app/mobius-cli
RUN --mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o mobiuscli ./cmd/mobiuscli

# Final stage
FROM alpine:3.20
RUN apk --no-cache add ca-certificates && \
	addgroup -S app && adduser -S -G app app
WORKDIR /app
COPY --from=builder /app/mobius-cli/mobiuscli /app/mobiuscli
USER app

# Run the CLI
ENTRYPOINT ["/app/mobiuscli"]
