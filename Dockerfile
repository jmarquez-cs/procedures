# Start with the official Golang base image
FROM golang:1.17-alpine as builder

# Set the working directory
WORKDIR /app

# Copy the Go project files into the container
COPY . .

# Build the binary
RUN go build -o sqlprocessor-cli ./cmd/sqlprocessor-cli/main.go

# Use a lightweight base image for the final stage
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/sqlprocessor-cli /app/sqlprocessor-cli

# Set the entrypoint for the container
ENTRYPOINT ["/app/sqlprocessor-cli"]

# By default, display the help message
CMD ["--help"]
