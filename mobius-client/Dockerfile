# Multi-stage build for Mobius Client
FROM golang:1.24.4-alpine AS builder

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Prepare cached deps
COPY go.work go.work.sum ./
COPY mobius-client/go.mod ./mobius-client/
COPY mobius-client/go.sum ./mobius-client/
COPY mobius-server/go.mod ./mobius-server/
COPY mobius-server/go.sum ./mobius-server/
COPY mobius-cli/go.mod ./mobius-cli/
COPY mobius-cli/go.sum ./mobius-cli/
COPY mobius-cocoon/go.mod ./mobius-cocoon/
COPY shared/go.mod ./shared/
COPY shared/go.sum ./shared/
RUN --mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	go work sync && \
	(cd mobius-client && go mod download) && \
	(cd mobius-server && go mod download) && \
	(cd mobius-cli && go mod download) && \
	(cd mobius-cocoon && go mod download) && \
	(cd shared && go mod download)

# Copy sources
COPY . .

# Build the client
WORKDIR /app/mobius-client
RUN --mount=type=cache,target=/go/pkg/mod \
	--mount=type=cache,target=/root/.cache/go-build \
	CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags "-s -w" -o client ./cmd/client

# Final stage
FROM alpine:3.20
RUN apk --no-cache add ca-certificates && \
	addgroup -S app && adduser -S -G app app
WORKDIR /app
COPY --from=builder /app/mobius-client/client /app/client
USER app

# Expose the client port if needed
EXPOSE 8081

# Run the client
ENTRYPOINT ["/app/client"]
