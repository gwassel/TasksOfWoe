# Use the official Golang image as the base image
FROM golang:1.23-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o task-tracker ./cmd/service

# Use a minimal Alpine image for the final stage
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/task-tracker .

# Copy the migration files
COPY migrations ./migrations

# Create the log directory and file
RUN mkdir -p /var/log/task-tracker && \
    touch /var/log/task-tracker/app.log /var/log/task-tracker/error.log && \
    chmod -R 777 /var/log/task-tracker

# Expose the port (if needed)
EXPOSE 8080

# Command to run the application
CMD ["./task-tracker"]