# Step 1: Build the Go binary
FROM golang:1.23.4-alpine AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum for dependency resolution
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the rest of the source code into the container
COPY . .

# Build the Go app
RUN go build -o main .

# Step 2: Create a smaller image to run the app
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the pre-built binary from the builder image
COPY --from=builder /app/main .

# Copy the .env file (if it exists)
# COPY .env .env

# Expose port 8080 (or whichever port your app uses)
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
