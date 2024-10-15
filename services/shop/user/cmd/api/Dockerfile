# Use the official Golang image to build the app
FROM golang:1.23 AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the container
COPY services/shop/user/go.mod services/shop/user/go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire code directory to the container
COPY services/shop/user .

# Build the Go app
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o user_service ./cmd/api

# Use a minimal base image to package the built app
FROM alpine:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the pre-built binary from the builder stage
COPY --from=builder /app/user_service .

# Copy configuration files from the builder stage
COPY services/shop/user/config/api_config.yaml /app/config/api_config.yaml

# Run the executable
CMD ["./user_service"]
