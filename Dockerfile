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

# Create a default environment file that can be overridden
RUN echo "ENVIRONMENT=development\n\
HOST=0.0.0.0\n\
PORT=8080\n\
DB_USERNAME=devuser\n\
DB_PASSWORD=devpass\n\
DB_HOSTNAME=postgres\n\
DB_PORT=5432\n\
DB_NAME=devdb\n\
JWT_SECRET=askimaskimaskimasecurelongersecret1234\n\
TOKEN_HOUR_LIFESPAN=24\n\
API_SECRET=jagiya\n\
VERSION=1\n\
GIN_MODE=release" > /app/dev.env

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["./gin-sqlc-app"]