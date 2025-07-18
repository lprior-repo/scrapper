# Multi-stage Docker build for GitHub Codeowners Visualization Tool

# Stage 1: Build Go application
FROM golang:1.24-alpine AS go-builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o overseer .

# Stage 2: Build React frontend
FROM node:20-alpine AS ui-builder

# Set working directory
WORKDIR /app/ui

# Copy package files
COPY ui/package.json ui/bun.lock* ./

# Install bun
RUN npm install -g bun@1.2.7

# Install dependencies
RUN bun install

# Copy UI source
COPY ui/ .

# Build the frontend
RUN bun run build

# Stage 3: Final runtime image
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 overseer && \
    adduser -D -s /bin/sh -u 1000 -G overseer overseer

# Create necessary directories
RUN mkdir -p /app/ui/dist && \
    chown -R overseer:overseer /app

# Set working directory
WORKDIR /app

# Copy built Go application
COPY --from=go-builder /app/overseer .

# Copy built frontend
COPY --from=ui-builder /app/ui/dist ./ui/dist/

# Copy UI server script
COPY --from=ui-builder /app/ui/server.ts ./ui/

# Set ownership
RUN chown -R overseer:overseer /app

# Switch to non-root user
USER overseer

# Expose ports
EXPOSE 8081 3000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8081/api/health || exit 1

# Default command
CMD ["./overseer", "api"]
