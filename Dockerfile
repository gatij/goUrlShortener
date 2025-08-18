FROM golang:1.22-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o urlshortener ./cmd/server

# Use a minimal alpine image for the final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy the binary from the builder stage
COPY --from=builder /app/urlshortener .

# Expose the application port
EXPOSE 3000

# Command to run the executable
CMD ["./urlshortener"]
