FROM golang:1.21-alpine AS builder

# Set working directory
WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app ./cmd/api

# Use a minimal alpine image for the final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata postgresql-client

# Set working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/app .

# Copy configuration files
COPY --from=builder /app/.env.example ./.env
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/scripts ./scripts

# Make scripts executable
RUN chmod +x ./scripts/run_migrations.sh

# Expose the application port
EXPOSE 8080

# Set the entrypoint
ENTRYPOINT ["./app"]
