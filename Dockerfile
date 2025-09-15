# Build stage
FROM golang:1.24.5-alpine AS builder

# Set working directory
WORKDIR /app

# Install dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o covo-bot main.go

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates tzdata

# Create app directory
WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/covo-bot .

# Copy JSON data files
COPY --from=builder /app/jsonfile ./jsonfile

# Create logs directory
RUN mkdir -p logs

# Set timezone
ENV TZ=Asia/Tehran

# Expose port (if needed)
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD pgrep covo-bot || exit 1

# Run the application
CMD ["./covo-bot"]
