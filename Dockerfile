# Stage 1: Build the Go executable
# Pin specific Go version for reproducible builds and security tracking
FROM golang:1.23-alpine AS builder

WORKDIR /server

# Copy dependency files first for better layer caching
COPY go.mod go.sum ./

# Download dependencies with verification
RUN go mod download && go mod verify

# Copy source code
COPY . .

# Accept version as build argument
ARG VERSION=dev

RUN CGO_ENABLED=0 go build -ldflags="-s -w -X github.com/aidanuno/qrkit/internal/server.Version=${VERSION}" -o qrkit

# Stage 2: Minimal runtime image
FROM scratch

COPY --from=builder /server/qrkit /server/qrkit

# EXPOSE is documentation for HTTP mode (use --port flag to enable)
# Default behavior without --port is STDIO mode (no network port needed)
EXPOSE 3000

# Set the working directory
WORKDIR /server

# Run as non-root user (nobody:nobody = 65534:65534)
# Note: scratch image doesn't have /etc/passwd, so we use numeric UID
USER 65534:65534

# Default: STDIO mode (docker run image)
# HTTP mode: docker run -p 3000:3000 image --port 3000
ENTRYPOINT ["./qrkit"]