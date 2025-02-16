FROM golang:1.23 AS builder

# Set the working directory
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application
RUN go build -o main .

# Use Alpine as the minimal base image for final execution
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Install necessary dependencies (if any)
RUN apk --no-cache add ca-certificates

# Expose port (change if needed)
EXPOSE 8080

# Command to run the application
CMD ["./main"]
