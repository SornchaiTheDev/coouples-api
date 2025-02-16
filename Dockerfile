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

# Use a minimal base image for final execution
FROM gcr.io/distroless/base-debian12

# Set the working directory
WORKDIR /root/

# Copy the built binary from the builder stage
COPY --from=builder /app/main .

# Expose port (change if needed)
EXPOSE 8080

# Command to run the application
CMD ["./main"]
