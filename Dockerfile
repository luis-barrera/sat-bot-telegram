# Use Golang as a build stage
FROM golang:1.21 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum separately for better caching
COPY go.mod go.sum ./

# Download dependencies only (cached if go.mod/go.sum didn't change)
RUN go mod download

# Copy the remaining project files
COPY . .

# Build the application
RUN go build -o /sat_telegram_bot

# Use a minimal runtime image
FROM debian:latest

# Set the working directory
WORKDIR /

# CACerts
RUN apt-get update && apt-get install -y ca-certificates
COPY certs /usr/local/share/ca-certificates/
RUN update-ca-certificates

# Copy the built binary from the builder stage
COPY --from=builder /sat_telegram_bot /sat_telegram_bot

# Run the Go application
ENTRYPOINT ["/sat_telegram_bot"]
