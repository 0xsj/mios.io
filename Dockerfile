# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev make git

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN make build

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/bin/gin-sqlc-app .
# Copy the environment file
COPY dev.env .

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./gin-sqlc-app"]