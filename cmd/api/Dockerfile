# Stage 1: Build the Go application
FROM golang:1.22 AS builder

# Install bash using apt-get since this is likely a Debian/Ubuntu-based image
RUN apt-get update && apt-get install -y bash

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Add wait-for-it.sh and make it executable
#COPY wait-for-it.sh /wait-for-it.sh
#RUN chmod +x /wait-for-it.sh

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo -o api ./cmd/api

# Stage 2: Create a minimal image with just the binary
FROM alpine:latest

WORKDIR /app

# Copy the built binary and wait-for-it.sh from the builder stage
COPY --from=builder /app .

# Expose the port that the app will run on
EXPOSE 8080

# Run the application using sh and wait-for-it.sh
CMD ["./api"]